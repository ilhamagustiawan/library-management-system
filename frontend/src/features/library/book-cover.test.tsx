import { fireEvent, render } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { BookCover } from "./book-cover";

describe("BookCover", () => {
  it("reveals a loaded catalog cover and keeps the generated cover as fallback", () => {
    const { container } = render(
      <BookCover
        coverUrl="https://covers.openlibrary.org/b/id/10521270-M.jpg"
        title="Clean Code"
      />,
    );
    const image = container.querySelector("img");
    if (image === null) throw new Error("Expected catalog cover image.");

    expect(image).toHaveAttribute("loading", "lazy");
    fireEvent.load(image);
    expect(image).toHaveAttribute("data-loaded", "true");

    fireEvent.error(image);
    expect(image).not.toBeVisible();
  });

  it("uses only the generated cover when no URL exists", () => {
    const { container } = render(<BookCover coverUrl={null} title="Clean Code" />);

    expect(container.querySelector("img")).toBeNull();
    expect(container).toHaveTextContent("Clean Code");
  });

  it("provides a large proportional cover for dashboard loan rows", () => {
    const { container } = render(
      <BookCover coverUrl={null} size="loan" title="The Alchemist" />,
    );

    expect(container.firstElementChild).toHaveAttribute("data-size", "loan");
    expect(container.firstElementChild).toHaveClass("aspect-[2/3]", "w-20", "sm:w-24");
    expect(container).toHaveTextContent("The Alchemist");
  });
});
