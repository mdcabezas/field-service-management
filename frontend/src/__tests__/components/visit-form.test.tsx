import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import React from "react";
import { VisitForm } from "@/components/forms/visit-form";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ back: vi.fn(), push: vi.fn() }),
}));
vi.mock("@/hooks/use-visits", () => ({
  useCreateVisit: () => ({ mutateAsync: vi.fn().mockResolvedValue({}), isSubmitting: false }),
}));
vi.mock("@/hooks/use-visit-types", () => ({
  useVisitTypes: () => ({ data: [{ id: "vt1", name: "Instalación", code: "INSTALL" }], isLoading: false }),
}));
vi.mock("sonner", () => ({ toast: { success: vi.fn(), error: vi.fn() } }));

function wrapper({ children }: { children: React.ReactNode }) {
  return React.createElement(QueryClientProvider, { client: new QueryClient({ defaultOptions: { queries: { retry: false } } }) }, children);
}

describe("VisitForm", () => {
  beforeEach(() => vi.clearAllMocks());

  it("renders all form fields", () => {
    render(React.createElement(VisitForm), { wrapper });
    expect(screen.getByLabelText(/fecha programada/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/prioridad/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/estado/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/tipo/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/destino de facturación/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/fuente/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/notas/i)).toBeInTheDocument();
  });

  it("shows validation error when scheduled_at is empty", async () => {
    const user = userEvent.setup();
    render(React.createElement(VisitForm), { wrapper });
    await user.click(screen.getByRole("button", { name: /crear/i }));
    await waitFor(() => {
      expect(screen.getByText("La fecha es requerida")).toBeInTheDocument();
    });
  });

  it("populates visit type dropdown from hook", () => {
    render(React.createElement(VisitForm), { wrapper });
    const typeSelect = screen.getByLabelText(/tipo/i);
    expect(typeSelect).toBeInTheDocument();
    expect(screen.getByText("Instalación (INSTALL)")).toBeInTheDocument();
  });

  it("renders create button", () => {
    render(React.createElement(VisitForm), { wrapper });
    expect(screen.getByRole("button", { name: /crear/i })).toBeInTheDocument();
  });

  it("renders cancel button", () => {
    render(React.createElement(VisitForm), { wrapper });
    expect(screen.getByRole("button", { name: /cancelar/i })).toBeInTheDocument();
  });

  it("has default values for priority and status", () => {
    render(React.createElement(VisitForm), { wrapper });
    const prioritySelect = screen.getByLabelText(/prioridad/i) as HTMLSelectElement;
    const statusSelect = screen.getByLabelText(/estado/i) as HTMLSelectElement;
    expect(prioritySelect.value).toBe("normal");
    expect(statusSelect.value).toBe("scheduled");
  });
});
