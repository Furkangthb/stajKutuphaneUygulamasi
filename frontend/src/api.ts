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

  const text = await res.text();

  if (!res.ok) {
    let message = res.statusText;
    if (text) {
      try {
        const parsed = JSON.parse(text);
        message = parsed.Hata || parsed.message || text;
      } catch {
        message = text;
      }
    }
    throw new Error(message);
  }

  if (!text) return undefined as T;

  try {
    return JSON.parse(text) as T;
  } catch {
    return text as unknown as T;
  }
}

export const api = {
  register: (data: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    phone: string;
  }) => req("POST", "/users/register", data),

  // Giriş işlemi
  login: (data: { email: string; password: string }) =>
    req<any>("POST", "/auth/login", data),

  logout: () => req("POST", "/auth/logout"),

  // Kullanıcılar
  getUsers: (page = 1, page_size = 20) =>
    req<UserListResponse>("GET", `/users?page=${page}&page_size=${page_size}`),
  deleteUser: (id: number) => req<void>("DELETE", `/users/${id}`),
  updateUser: (id: number, data: Partial<User>) => req<User>("PUT", `/users/${id}`, data),
  getUser: (id: number) => req<User>("GET", `/users/${id}`),
  changePassword: (id: number, data: { current_password: string; new_password: string }) =>
    req<{ mesaj: string }>("PUT", `/users/${id}/password`, data),

  // Kitaplar
  getBooks: (limit: number = 100) => req<Book[]>("GET", `/books?page=1&page_size=${limit}`),
  getBook: (id: number) => req<Book>("GET", `/books/${id}`),
  addBook: (data: Omit<Book, "id">) => req<Book>("POST", "/books", data),
  updateBook: (id: number, data: Partial<Book>) => req<Book>("PUT", `/books/${id}`, data),
  deleteBook: (id: number) => req("DELETE", `/books/${id}`, undefined),

  // Rezervasyonlar
  createReservation: (data: { book_id: number; user_id?: number }) =>
    req<Reservation>("POST", "/reservation", data),
  updateReservation: (id: number, data: Partial<Reservation>) =>
    req<Reservation>("PUT", `/reservation/${id}`, data),
  getUserReservations: (userId: number) =>
    req<Reservation[]>("GET", `/reservation/${userId}`),

  getAllReservations: () => req<ReservationFull[]>("GET", "/reservations"),

  // Sohbet
  getChatHistory: (limit = 50) =>
    req<{ data: ChatMessage[] }>("GET", `/chat/history?limit=${limit}`),
};


export interface User {
  id: number;
  first_name: string;
  last_name: string;
  phone?: string;
  email: string;
  role: "user" | "admin";
  created_at?: string;
}

export interface UserListResponse {
  data: User[];
  page: number;
}

export interface Book {
  id: number;
  isbn: string;
  title: string;
  author: string;
  genre: string;
  publish_date: string;
  description: string;
  available?: boolean;
}

export interface Reservation {
  id: number;
  user_id: number;
  book_id: number;
  status: "active" | "completed" | "cancelled" | "expired";
  reserved_at?: string;
  due_date?: string;
  book?: Book;
  user?: User;
}

export interface ReservationFull {
  id: number;
  user_id: number;
  book_id: number;
  status: "active" | "completed" | "cancelled" | "expired";
  reserved_at: string;
  due_date: string;
  first_name: string;
  last_name: string;
  book_title: string;
  book_author: string;
}

export interface ChatMessage {
  ID: number;
  UserID: number;
  Role: "user" | "assistant";
  Message: string;
  CreatedAt: string;
}