import { z } from "zod";

const password = z
  .string()
  .min(8, "Password must contain at least 8 characters.");

const loginSchema = z.object({
  email: z.email("Enter a valid email address."),
  password,
});

export type LoginInput = z.infer<typeof loginSchema>;

export const LoginInput = {
  schema: loginSchema,
} as const;

const registerSchema = z
  .object({
    name: z.string().trim().min(2, "Name must contain at least 2 characters."),
    email: z.email("Enter a valid email address."),
    password,
    confirmPassword: z.string(),
    acceptsTerms: z
      .boolean()
      .refine((accepted) => accepted, "Accept the terms to create an account."),
  })
  .superRefine(({ confirmPassword, password }, context) => {
    if (confirmPassword !== password) {
      context.addIssue({
        code: "custom",
        message: "Passwords must match.",
        path: ["confirmPassword"],
      });
    }
  });

export type RegisterInput = z.infer<typeof registerSchema>;

export const RegisterInput = {
  schema: registerSchema,
} as const;
