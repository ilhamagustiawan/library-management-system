export function FormError({ id, message }: { id: string; message: string | undefined }) {
  if (message === undefined) return null;

  return (
    <p id={id} className="text-sm text-destructive">
      {message}
    </p>
  );
}
