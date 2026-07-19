import { z } from "zod";

const STORAGE_KEY = "libry.mock-session";

const schema = z.object({
  id: z.string().min(1),
  name: z.string().min(1),
  email: z.email(),
});

type SessionStorageReader = Pick<Storage, "getItem">;
type SessionStorageWriter = Pick<Storage, "removeItem" | "setItem">;

export type AuthSession = z.infer<typeof schema>;

let cachedSerialized: string | null | undefined;
let cachedSession: AuthSession | null = null;

export const AuthSession = {
  clear(storage: SessionStorageWriter) {
    storage.removeItem(STORAGE_KEY);
  },
  read(storage: SessionStorageReader): AuthSession | null {
    const stored = storage.getItem(STORAGE_KEY);

    if (stored === cachedSerialized) return cachedSession;

    cachedSerialized = stored;

    if (stored === null) {
      cachedSession = null;
      return null;
    }

    try {
      const parsed: unknown = JSON.parse(stored);
      const result = schema.safeParse(parsed);
      cachedSession = result.success ? result.data : null;
      return cachedSession;
    } catch {
      cachedSession = null;
      return null;
    }
  },
  write(storage: SessionStorageWriter, session: AuthSession) {
    storage.setItem(STORAGE_KEY, JSON.stringify(session));
  },
} as const;
