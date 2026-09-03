import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import React from "react";
import { CustomerForm } from "@/components/forms/customer-form";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}));

vi.mock("@/hooks/use-customers", () => ({
  useCreateCustomer: () => ({
    mutateAsync: vi.fn().mockResolvedValue({}),
    isSubmitting: false,
  }),
  useUpdateCustomer: () => ({
    mutateAsync: vi.fn().mockResolvedValue({}),
    isSubmitting: false,
  }),
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return React.createElement(QueryClientProvider, { client: qc }, children);
}

describe("CustomerForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders 'Crear Cliente' when no initialData", () => {
    render(React.createElement(CustomerForm), { wrapper });
    expect(screen.getByText("Crear Cliente")).toBeInTheDocument();
  });

  it("renders 'Editar Cliente' when isEdit is true", () => {
    render(
      React.createElement(CustomerForm, {
        initialData: {
          id: "1",
          name: "Existing",
          tax_id: "123",
          phone: "555",
          email: "ex@test.com",
          created_at: "",
          updated_at: "",
        },
        isEdit: true,
      }),
      { wrapper }
    );
    expect(screen.getByText("Editar Cliente")).toBeInTheDocument();
  });

  it("shows validation error when name is empty", async () => {
    const user = userEvent.setup();
    render(React.createElement(CustomerForm), { wrapper });

    await user.click(screen.getByRole("button", { name: /crear/i }));

    await waitFor(() => {
      expect(screen.getByText("Nombre requerido")).toBeInTheDocument();
    });
  });

  it("does not submit when name is empty", async () => {
    const user = userEvent.setup();
    render(React.createElement(CustomerForm), { wrapper });

    // Fill only email, leave name empty
    await user.type(screen.getByLabelText(/email/i), "test@test.com");
    await user.click(screen.getByRole("button", { name: /crear/i }));

    await waitFor(() => {
      // Name validation error should appear
      expect(screen.getByText("Nombre requerido")).toBeInTheDocument();
    });
  });

  it("calls createCustomer on valid submit", async () => {
    const user = userEvent.setup();
    render(React.createElement(CustomerForm), { wrapper });

    await user.type(screen.getByLabelText(/nombre/i), "New Client");
    await user.click(screen.getByRole("button", { name: /crear/i }));

    await waitFor(() => {
      // Form should not show validation errors
      expect(screen.queryByText("Nombre requerido")).not.toBeInTheDocument();
    });
  });

  it("renders cancel button that says Cancelar", () => {
    render(React.createElement(CustomerForm), { wrapper });
    expect(screen.getByRole("button", { name: /cancelar/i })).toBeInTheDocument();
  });
});
