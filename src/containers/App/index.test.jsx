import React from "react";
import { render } from "@testing-library/react";
import App from "./index";

test("renders learn react link", () => {
  const { getByText } = render(<App />);
  const linkElement = getByText(/herodote/i);
  expect(linkElement).toBeInTheDocument();
});
