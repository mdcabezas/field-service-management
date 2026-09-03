import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import React from "react";
import { UserForm } from "@/components/forms/user-form";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));
vi.mock("@/hooks/use-users", () => ({
  useCreateUser: () => ({ mutateAsync: vi.fn().mockResolvedValue({}), isSubmitting: false }),
  useUpdateUser: () => ({ mutateAsync: vi.fn().mockResolvedValue({}), isSubmitting: false }),
}));
vi.mock("sonner", () => ({ toast: { success: vi.fn(), error: vi.fn() } }));

function wrapper({ children }: { children: React.ReactNode }) {
  return React.createElement(QueryClientProvider, { client: new QueryClient({ defaultOptions: { queries: { retry: false } } }) }, children);
}

describe("UserForm", () => {
  beforeEach(() => vi.clearAllMocks());

  it("renders 'Crear Usuario' when no initialData", () => {
    render(React.createElement(UserForm), { wrapper });
    expect(screen.getByText("Crear Usuario")).toBeInTheDocument();
  });

  it("renders 'Editar Usuario' in edit mode", () => {
    render(React.createElement(UserForm, { initialData: { id: "1", email: "a@b.com", name: "A", role: "admin", created_at: "" }, isEdit: true }), { wrapper });
    expect(screen.getByText("Editar Usuario")).toBeInTheDocument();
  });

  it("shows email validation error when empty", async () => {
    const user = userEvent.setup();
    render(React.createElement(UserForm), { wrapper });
    await user.click(screen.getByRole("button", { name: /crear/i }));
    await waitFor(() => {
      expect(screen.getByText("Email inválido")).toBeInTheDocument();
    });
  });

  it("shows name validation error when empty", async () => {
    const user = userEvent.setup();
    render(React.createElement(UserForm), { wrapper });
    await user.type(screen.getByLabelText(/email/i), "a@b.com");
    await user.click(screen.getByRole("button", { name: /crear/i }));
    await waitFor(() => {
      expect(screen.getByText("Nombre requerido")).toBeInTheDocument();
    });
  });

  it("does not show password field in edit mode", () => {
    render(React.createElement(UserForm, { initialData: { id: "1", email: "a@b.com", name: "A", role: "admin", created_at: "" }, isEdit: true }), { wrapper });
    expect(screen.queryByLabelText(/contraseña/i)).not.toBeInTheDocument();
  });

  it("shows password field in create mode", () => {
    render(React.createElement(UserForm), { wrapper });
    expect(screen.getByLabelText(/contraseña/i)).toBeInTheDocument();
  });

  it("renders cancel button", () => {
    render(React.createElement(UserForm), { wrapper });
    expect(screen.getByRole("button", { name: /cancelar/i })).toBeInTheDocument();
  });
});
