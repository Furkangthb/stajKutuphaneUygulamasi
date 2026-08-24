const BASE = "/api";

async function req<T>(method: string, path: string, body?: unknown): Promise<T> {
  const token = localStorage.getItem("token");
  const res = await fetch(`${BASE}${path}`, {
    method,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    ...(body ? { body: JSON.stringify(body) } : {}),
  });
  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || res.statusText);
  }
  const text = await res.text();
  return (text ? JSON.parse(text) : undefined) as T;
}

export const api = {
  register: (data: { username: string; email: string; password: string }) =>
    req("POST", "/users/register", data),

  login: (data: { email: string; password: string }) =>
    req<{ token: string; user: User }>("POST", "/auth/login", data),

  logout: () => req("POST", "/auth/logout"),

  // Users (admin)
  getUsers: () => req<User[]>("GET", "/users"),
  deleteUser: (id: number) => req("DELETE", `/users/${id}`),
  updateUser: (id: number, data: Partial<User>) => req("PUT", `/users/${id}`, data),
  getUser: (id: number) => req<User>("GET", `/users/${id}`),

  // Books
  getBooks: () => req<Book[]>("GET", "/books"),
  getBook: (id: number) => req<Book>("GET", `/books/${id}`),
  addBook: (data: Omit<Book, "id">) => req<Book>("POST", "/books", data),
  updateBook: (id: number, data: Partial<Book>) => req<Book>("PUT", `/books/${id}`, data),
  deleteBook: (id: number) => req("DELETE", `/books/${id}`, undefined),

  // Reservations
  createReservation: (data: { book_id: number }) =>
    req<Reservation>("POST", "/reservation", data),
  updateReservation: (id: number, data: Partial<Reservation>) =>
    req<Reservation>("PUT", `/reservation/${id}`, data),
  getUserReservations: (userId: number) =>
    req<Reservation[]>("GET", `/reservation/${userId}`),
};

// Tipleri ayrı export ediyoruz
export interface User {
  id: number;
  username: string;
  email: string;
  role: "user" | "admin";
  created_at?: string;
}

export interface Book {
  id: number;
  title: string;
  author: string;
  isbn?: string;
  genre?: string;
  published_year?: number;
  available: boolean;
  description?: string;
}

export interface Reservation {
  id: number;
  user_id: number;
  book_id: number;
  status: "pending" | "active" | "returned" | "cancelled";
  reserved_at?: string;
  returned_at?: string;
  book?: Book;
  user?: User;
}