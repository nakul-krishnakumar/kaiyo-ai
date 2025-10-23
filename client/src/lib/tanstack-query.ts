const API_URL: string = (import.meta as any).env?.VITE_API_URL ?? "";

export function getTokens() {
  return {
    access: localStorage.getItem("access_token"),
    accessExpiresAt: Number(localStorage.getItem("access_expires_at")),
  };
}

export function setTokens(data: any) {
  if (data?.access_token) {
    localStorage.setItem("access_token", data.access_token);
  }

  if (data?.expires_in) {
    const expiresAt = Date.now() + Number(data.expires_in) * 1000;
    localStorage.setItem("access_expires_at", String(expiresAt));
  }
}

export function isAuthenticated(): boolean {
  const { access, accessExpiresAt } = getTokens();
  if (!access) return false;
  if (accessExpiresAt && Date.now() > accessExpiresAt) return false; // token expired
  return true;
}

export async function logout() {
  try {
    if (API_URL) {
      await fetch(`${API_URL}/auth/logout`, {
        method: "GET",
        credentials: "include",
      });
    }
  } catch {
    // ignore network errors
  }
  localStorage.removeItem("access_token");
  localStorage.removeItem("access_expires_at");
  window.location.href = "/login";
}

async function refreshAccessToken() {
  if (!API_URL) return null;
  try {
    const res = await fetch(`${API_URL}/auth/refresh`, {
      method: "GET",
      credentials: "include",
    });
    if (!res.ok) return null;
    const data = await res.json();
    setTokens(data);
    return data;
  } catch {
    return null;
  }
}

export async function fetchClient(input: string, init?: RequestInit): Promise<Response> {
  const { access } = getTokens();
  const base = API_URL;

  const headers = new Headers(init?.headers);
  if (access) headers.set("Authorization", `Bearer ${access}`);

  let response = await fetch(`${base}${input}`, {
    ...init,
    headers,
    credentials: "include",
  });

  // Retry if unauthorized
  if (response.status === 401) {
    const refreshed = await refreshAccessToken();
    if (refreshed?.access_token) {
      const retryHeaders = new Headers(init?.headers);
      retryHeaders.set("Authorization", `Bearer ${refreshed.access_token}`);
      response = await fetch(`${base}${input}`, {
        ...init,
        headers: retryHeaders,
        credentials: "include",
      });
    } else {
      await logout();
      throw new Error("Session expired. Please log in again.");
    }
  }

  return response;
}
