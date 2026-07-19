import { createCipheriv, createDecipheriv, createHash, randomBytes } from "node:crypto";

import { z } from "zod";

export type OpenSealedValueResult<T> =
  | { status: "valid"; value: T }
  | { status: "invalid" };

function encryptionKey(secret: string) {
  if (secret.length < 32) {
    throw new Error("AUTH_SESSION_SECRET must contain at least 32 characters");
  }
  return createHash("sha256").update(secret, "utf8").digest();
}

function seal<T>(value: T, schema: z.ZodType<T>, secret: string) {
  const validated = schema.parse(value);
  const iv = randomBytes(12);
  const cipher = createCipheriv("aes-256-gcm", encryptionKey(secret), iv);
  const ciphertext = Buffer.concat([
    cipher.update(JSON.stringify(validated), "utf8"),
    cipher.final(),
  ]);
  const tag = cipher.getAuthTag();
  return ["v1", iv.toString("base64url"), ciphertext.toString("base64url"), tag.toString("base64url")].join(".");
}

function open<T>(value: string, schema: z.ZodType<T>, secret: string): OpenSealedValueResult<T> {
  const parts = value.split(".");
  const version = parts[0];
  const encodedIV = parts[1];
  const encodedCiphertext = parts[2];
  const encodedTag = parts[3];
  if (
    parts.length !== 4 ||
    version !== "v1" ||
    encodedIV === undefined ||
    encodedCiphertext === undefined ||
    encodedTag === undefined
  ) {
    return { status: "invalid" };
  }

  try {
    const decipher = createDecipheriv(
      "aes-256-gcm",
      encryptionKey(secret),
      Buffer.from(encodedIV, "base64url"),
    );
    decipher.setAuthTag(Buffer.from(encodedTag, "base64url"));
    const plaintext = Buffer.concat([
      decipher.update(Buffer.from(encodedCiphertext, "base64url")),
      decipher.final(),
    ]).toString("utf8");
    const parsed: unknown = JSON.parse(plaintext);
    const result = schema.safeParse(parsed);
    return result.success ? { status: "valid", value: result.data } : { status: "invalid" };
  } catch {
    return { status: "invalid" };
  }
}

export const SealedValue = { open, seal } as const;
