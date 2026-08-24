import { useEffect, useState } from "react";
import { api, type Reservation } from "../api";

const STATUS_LABELS: Record<string, string> = {
  pending: "Bekliyor",
  active: "Aktif",
  returned: "İade Edildi",
  cancelled: "İptal",
};

const STATUS_COLORS: Record<string, { bg: string; text: string }> = {
  pending: { bg: "#FEF3C7", text: "#92400E" },
  active: { bg: "#D1FAE5", text: "#065F46" },
  returned: { bg: "#DBEAFE", text: "#1E40AF" },
  cancelled: { bg: "#FEE2E2", text: "#991B1B" },
};

const STATUS_FILTERS = ["Tümü", "Bekliyor", "Aktif", "İade Edildi", "İptal"];
const STATUS_MAP: Record<string, string> = {
  "Bekliyor": "pending", "Aktif": "active", "İade Edildi": "returned", "İptal": "cancelled",
};

export default function AllReservationsPage() {
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [statusFilter, setStatusFilter] = useState("Tümü");
  const [search, setSearch] = useState("");
  const [updating, setUpdating] = useState<number | null>(null);
  const [toast, setToast] = useState("");

  useEffect(() => {
    // Admin: fetch all reservations via a listing endpoint
    // Backend route: GET /api/reservation/:id (user-based)
    // For admin "all reservations" we use a conventional /api/reservations endpoint
    fetch("/api/reservations", {
      headers: { Authorization: `Bearer ${localStorage.getItem("token") ?? ""}` },
    })
      .then(async (res) => {
        if (!res.ok) throw new Error(await res.text());
        return res.json();
      })
      .then(setReservations)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  const showToast = (msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(""), 3000);
  };

  const updateStatus = async (id: number, status: Reservation["status"]) => {
    setUpdating(id);
    try {
      const updated = await api.updateReservation(id, { status });
      setReservations((rs) => rs.map((r) => (r.id === id ? { ...r, ...updated } : r)));
      showToast("Rezervasyon güncellendi");
    } catch (e: unknown) {
      showToast(e instanceof Error ? e.message : "Güncelleme başarısız");
    } finally {
      setUpdating(null);
    }
  };

  const filtered = reservations.filter((r) => {
    const matchStatus = statusFilter === "Tümü" || r.status === STATUS_MAP[statusFilter];
    const matchSearch =
      search === "" ||
      r.book?.title?.toLowerCase().includes(search.toLowerCase()) ||
      r.user?.username?.toLowerCase().includes(search.toLowerCase());
    return matchStatus && matchSearch;
  });

  const counts = {
    total: reservations.length,
    active: reservations.filter((r) => r.status === "active").length,
    pending: reservations.filter((r) => r.status === "pending").length,
    returned: reservations.filter((r) => r.status === "returned").length,
  };

  return (
    <div className="h-full flex flex-col">
      <div className="px-8 pt-8 pb-6">
        <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
          Tüm Rezervasyonlar
        </h1>
        <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
          {reservations.length} rezervasyon
        </p>
      </div>

      {/* Summary cards */}
      <div className="px-8 pb-5 grid grid-cols-2 sm:grid-cols-4 gap-3">
        {[
          { label: "Toplam", value: counts.total, bg: "#F8F7F3", text: "var(--foreground)" },
          { label: "Aktif", value: counts.active, bg: "#D1FAE5", text: "#065F46" },
          { label: "Bekleyen", value: counts.pending, bg: "#FEF3C7", text: "#92400E" },
          { label: "İade", value: counts.returned, bg: "#DBEAFE", text: "#1E40AF" },
        ].map(({ label, value, bg, text }) => (
          <div key={label} className="rounded-xl px-5 py-4" style={{ backgroundColor: bg, border: "1px solid var(--border)" }}>
            <p className="font-display text-2xl" style={{ color: text }}>{value}</p>
            <p className="text-xs font-semibold mt-0.5" style={{ color: text, opacity: 0.65 }}>{label}</p>
          </div>
        ))}
      </div>

      {/* Filters */}
      <div className="px-8 pb-4 flex flex-wrap gap-3">
        <div className="relative flex-1 min-w-48">
          <svg className="absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="var(--muted-foreground)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="11" cy="11" r="8" /><line x1="21" y1="21" x2="16.65" y2="16.65" />
          </svg>
          <input
            type="text"
            placeholder="Kitap veya kullanıcı ara…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full pl-9 pr-4 py-2.5 rounded-xl text-sm outline-none"
            style={{ backgroundColor: "var(--card)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
          />
        </div>
        <div className="flex gap-1.5 flex-wrap">
          {STATUS_FILTERS.map((s) => (
            <button
              key={s}
              onClick={() => setStatusFilter(s)}
              className="px-3 py-2 rounded-lg text-xs font-semibold transition-all cursor-pointer"
              style={{
                backgroundColor: statusFilter === s ? "var(--primary)" : "var(--card)",
                color: statusFilter === s ? "var(--primary-foreground)" : "var(--muted-foreground)",
                border: `1.5px solid ${statusFilter === s ? "var(--primary)" : "var(--border)"}`,
              }}
            >
              {s}
            </button>
          ))}
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-8 pb-8">
        {loading ? (
          <LoadingSkeleton />
        ) : error ? (
          <div className="flex flex-col items-center justify-center h-48 gap-2">
            <p className="text-sm font-semibold" style={{ color: "#DC2626" }}>Veri yüklenemedi</p>
            <p className="text-xs" style={{ color: "var(--muted-foreground)" }}>{error}</p>
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex items-center justify-center h-48">
            <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>Rezervasyon bulunamadı</p>
          </div>
        ) : (
          <div className="rounded-2xl overflow-hidden" style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}>
            <table className="w-full text-sm">
              <thead>
                <tr style={{ borderBottom: "1px solid var(--border)" }}>
                  {["#", "Kitap", "Kullanıcı", "Durum", "Rezervasyon Tarihi", ""].map((h) => (
                    <th key={h} className="px-5 py-3 text-left font-semibold text-xs uppercase tracking-wider" style={{ color: "var(--muted-foreground)" }}>
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {filtered.map((r, i) => {
                  const sc = STATUS_COLORS[r.status] ?? { bg: "#F3F4F6", text: "#6B7280" };
                  return (
                    <tr
                      key={r.id}
                      style={{ borderBottom: i < filtered.length - 1 ? "1px solid var(--border)" : "none" }}
                      className="transition-colors duration-100 hover:bg-[var(--muted)]"
                    >
                      <td className="px-5 py-3 font-mono text-xs" style={{ color: "var(--muted-foreground)" }}>#{r.id}</td>
                      <td className="px-5 py-3 font-semibold" style={{ color: "var(--foreground)" }}>
                        {r.book?.title ?? `Kitap #${r.book_id}`}
                      </td>
                      <td className="px-5 py-3">
                        <div className="flex items-center gap-2">
                          <div
                            className="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold uppercase shrink-0"
                            style={{ backgroundColor: "var(--secondary)", color: "var(--secondary-foreground)" }}
                          >
                            {r.user?.username?.[0] ?? "?"}
                          </div>
                          <span style={{ color: "var(--muted-foreground)" }}>
                            {r.user?.username ?? `Kullanıcı #${r.user_id}`}
                          </span>
                        </div>
                      </td>
                      <td className="px-5 py-3">
                        <span className="px-2 py-0.5 rounded-full text-xs font-bold" style={{ backgroundColor: sc.bg, color: sc.text }}>
                          {STATUS_LABELS[r.status]}
                        </span>
                      </td>
                      <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>
                        {r.reserved_at ? new Date(r.reserved_at).toLocaleDateString("tr-TR") : "—"}
                      </td>
                      <td className="px-5 py-3">
                        <select
                          value={r.status}
                          disabled={updating === r.id}
                          onChange={(e) => updateStatus(r.id, e.target.value as Reservation["status"])}
                          className="px-2 py-1.5 rounded-lg text-xs font-semibold outline-none cursor-pointer"
                          style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
                        >
                          <option value="pending">Bekliyor</option>
                          <option value="active">Aktif</option>
                          <option value="returned">İade Edildi</option>
                          <option value="cancelled">İptal</option>
                        </select>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {toast && (
        <div className="fixed bottom-6 right-6 px-5 py-3 rounded-xl text-sm font-semibold shadow-lg z-50" style={{ backgroundColor: "var(--sidebar)", color: "var(--sidebar-foreground)" }}>
          {toast}
        </div>
      )}
    </div>
  );
}

function LoadingSkeleton() {
  return (
    <div className="rounded-2xl overflow-hidden animate-pulse" style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}>
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="px-5 py-4 flex gap-4" style={{ borderBottom: "1px solid var(--border)" }}>
          <div className="h-4 rounded w-8" style={{ backgroundColor: "var(--muted)" }} />
          <div className="h-4 rounded flex-1" style={{ backgroundColor: "var(--muted)" }} />
          <div className="h-4 rounded w-24" style={{ backgroundColor: "var(--muted)" }} />
          <div className="h-4 rounded w-16" style={{ backgroundColor: "var(--muted)" }} />
        </div>
      ))}
    </div>
  );
}
