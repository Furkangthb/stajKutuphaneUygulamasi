import { useState, useEffect } from "react";
import type { User } from "./api";
import AuthPage from "./pages/AuthPage";
import BooksPage from "./pages/BooksPage";
import ReservationsPage from "./pages/ReservationsPage";
import ManageBooksPage from "./pages/ManageBooksPage";
import ManageUsersPage from "./pages/ManageUsersPage";
import AllReservationsPage from "./pages/AllReservationsPage";
import ChatPage from "./pages/ChatPage";
import Sidebar from "./components/Sidebar";

type Page = "books" | "reservations" | "manage-books" | "manage-users" | "all-reservations" | "chat";

export default function App() {
  const [user, setUser] = useState<User | null>(null);
  const [page, setPage] = useState<Page>("books");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("token");
    const storedUser = localStorage.getItem("user");

    // Eğer token ve kullanıcı bilgisi varsa oturumu açık tutuyoruz
    if (token && storedUser && storedUser !== "undefined") {
      try {
        setUser(JSON.parse(storedUser));
      } catch (err) {
        console.error("Kullanıcı bilgisi okunamadı", err);
      }
    }
    setLoading(false);
  }, []);

  const handleLogin = (u: User, token: string) => {
    localStorage.setItem("token", token); // Token'ı backend'e istek atarken kullanmak için kaydediyoruz
    localStorage.setItem("user", JSON.stringify(u));
    setUser(u);
    setPage("books");
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setUser(null);
  };

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center text-white" style={{ backgroundColor: "var(--background)" }}>
        Yükleniyor...
      </div>
    );
  }

  if (!user) {
    return <AuthPage onLogin={handleLogin} />;
  }

  const renderPage = () => {
    switch (page) {
      case "books":
        return <BooksPage userId={user.id} />;
      case "reservations":
        return <ReservationsPage userId={user.id} />;
      case "manage-books":
        return user.role === "admin" ? <ManageBooksPage /> : <BooksPage userId={user.id} />;
      case "manage-users":
        return user.role === "admin" ? <ManageUsersPage /> : <BooksPage userId={user.id} />;
      case "all-reservations":
        return user.role === "admin" ? <AllReservationsPage /> : <ReservationsPage userId={user.id} />;
      case "chat":
        return <ChatPage />;
      default:
        return <BooksPage userId={user.id} />;
    }
  };

  return (
    <div className="flex h-screen overflow-hidden" style={{ backgroundColor: "var(--background)" }}>
      <Sidebar
        user={user}
        active={page}
        onNavigate={(p: string) => setPage(p as Page)}
        onLogout={handleLogout}
      />
      <main className="flex-1 overflow-hidden" style={{ backgroundColor: "var(--background)" }}>
        {renderPage()}
      </main>
    </div>
  );
}