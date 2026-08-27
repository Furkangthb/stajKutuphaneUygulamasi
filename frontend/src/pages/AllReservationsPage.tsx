import { useEffect, useState } from "react";
import { api, type ReservationFull, type Reservation } from "../api";

const STATUS_LABELS: Record<string, string> = {
  pending: "Onay Bekliyor",
  active: "Aktif",
  completed: "Tamamlandı",
  cancelled: "İptal Edildi",
  expired: "Süresi Doldu",
};

const STATUS_COLORS: Record<string, { bg: string; text: string }> = {
  pending: { bg: "#EDE9FE", text: "#5B21B6" },
  active: { bg: "#D1FAE5", text: "#065F46" },
  completed: { bg: "#DBEAFE", text: "#1E40AF" },
  cancelled: { bg: "#FEE2E2", text: "#991B1B" },
  expired: { bg: "#FEF3C7", text: "#92400E" },
};

const STATUS_FILTERS = ["Tümü", "Onay Bekliyor", "Aktif", "Tamamlandı", "İptal Edildi", "Süresi Doldu"];
const STATUS_MAP: Record<string, string> = {
  "Onay Bekliyor": "pending", "Aktif": "active", "Tamamlandı": "completed", "İptal Edildi": "cancelled", "Süresi Doldu": "expired",
};

export default function AllReservationsPage() {
  const [reservations, setReservations] = useState<ReservationFull[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [statusFilter, setStatusFilter] = useState("Tümü");
  const [search, setSearch] = useState("");
  const [updating, setUpdating] = useState<number | null>(null);
  const [toast, setToast] = useState("");

  useEffect(() => {
    api.getAllReservations()
      .then((res: any) => {
        const dataList = Array.isArray(res) ? res : (res?.data || res?.reservations || []);
        setReservations(dataList || []);
      })
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
      await api.updateReservation(id, { status });
      setReservations((rs) => (rs || []).map((r) => (r?.id === id ? { ...r, status } : r)));
      showToast(status === "active" ? "Rezervasyon onaylandı" : status === "cancelled" ? "Rezervasyon reddedildi/iptal edildi" : "Rezervasyon güncellendi");
    } catch (e: unknown) {
      showToast(e instanceof Error ? e.message : "Güncelleme başarısız");
    } finally {
      setUpdating(null);
    }
  };

  const safeReservations = Array.isArray(reservations) ? reservations : [];

  const filtered = safeReservations.filter((r) => {
    if (!r) return false;
    const matchStatus = statusFilter === "Tümü" || r.status === STATUS_MAP[statusFilter];

    const bookTitle = r.book_title || "";
    const userName = `${r.first_name || ""} ${r.last_name || ""}`.trim();
    const searchTerm = search.toLowerCase();

    const matchSearch =
      search === "" ||
      bookTitle.toLowerCase().includes(searchTerm) ||
      userName.toLowerCase().includes(searchTerm);

    return matchStatus && matchSearch;
  });

  const counts = {
    total: safeReservations.length,
    pending: safeReservations.filter((r) => r?.status === "pending").length,
    active: safeReservations.filter((r) => r?.status === "active").length,
    completed: safeReservations.filter((r) => r?.status === "completed").length,
    cancelled: safeReservations.filter((r) => r?.status === "cancelled").length,
  };

  return (
    <div className="h-full flex flex-col">
      <div className="px-8 pt-8 pb-6">
        <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
          Tüm Rezervasyonlar
        </h1>
        <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
          {safeReservations.length} rezervasyon
        </p>
      </div>

      <div className="px-8 pb-5 grid grid-cols-2 sm:grid-cols-5 gap-3">
        {[
          { label: "Toplam", value: counts.total, bg: "#F8F7F3", text: "var(--foreground)" },
          { label: "Onay Bekliyor", value: counts.pending, bg: "#EDE9FE", text: "#5B21B6" },
          { label: "Aktif", value: counts.active, bg: "#D1FAE5", text: "#065F46" },
          { label: "Tamamlandı", value: counts.completed, bg: "#DBEAFE", text: "#1E40AF" },
          { label: "İptal Edilen", value: counts.cancelled, bg: "#FEE2E2", text: "#991B1B" },
        ].map(({ label, value, bg, text }) => (
          <div key={label} className="rounded-xl px-5 py-4" style={{ backgroundColor: bg, border: "1px solid var(--border)" }}>
            <p className="font-display text-2xl" style={{ color: text }}>{value}</p>
            <p className="text-xs font-semibold mt-0.5" style={{ color: text, opacity: 0.65 }}>{label}</p>
          </div>
        ))}
      </div>

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
                  {["#", "Kitap", "Kullanıcı", "Durum", "İşlem Tarihi", "İade Tarihi", ""].map((h) => (
                    <th key={h} className="px-5 py-3 text-left font-semibold text-xs uppercase tracking-wider" style={{ color: "var(--muted-foreground)" }}>
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {filtered.map((r, i) => {
                  if (!r) return null;
                  const sc = STATUS_COLORS[r.status] ?? { bg: "#F3F4F6", text: "#6B7280" };
                  const fullName = `${r.first_name || ""} ${r.last_name || ""}`.trim() || `Kullanıcı #${r.user_id}`;
                  const isOverdue = r.due_date && new Date(r.due_date) < new Date() && r.status === "active";
                  const isPending = r.status === "pending";

                  return (
                    <tr
                      key={r.id || i}
                      style={{ borderBottom: i < filtered.length - 1 ? "1px solid var(--border)" : "none" }}
                      className="transition-colors duration-100 hover:bg-[var(--muted)]"
                    >
                      <td className="px-5 py-3 font-mono text-xs" style={{ color: "var(--muted-foreground)" }}>#{r.id}</td>
                      <td className="px-5 py-3 font-semibold" style={{ color: "var(--foreground)" }}>
                        {r.book_title || `Kitap #${r.book_id}`}
                      </td>
                      <td className="px-5 py-3">
                        <div className="flex items-center gap-2">
                          <div
                            className="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold uppercase shrink-0"
                            style={{ backgroundColor: "var(--secondary)", color: "var(--secondary-foreground)" }}
                          >
                            {r.first_name ? r.first_name[0] : "?"}
                          </div>
                          <span style={{ color: "var(--muted-foreground)" }}>
                            {fullName}
                          </span>
                        </div>
                      </td>
                      <td className="px-5 py-3">
                        <span className="px-2 py-0.5 rounded-full text-xs font-bold" style={{ backgroundColor: sc.bg, color: sc.text }}>
                          {STATUS_LABELS[r.status] || "Bilinmiyor"}
                        </span>
                      </td>
                      <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>
                        {r.reserved_at ? new Date(r.reserved_at).toLocaleDateString("tr-TR") : "—"}
                      </td>
                      <td
                        className="px-5 py-3"
                        style={{
                          color: isOverdue ? "#DC2626" : "var(--muted-foreground)",
                          fontWeight: isOverdue ? "bold" : "500"
                        }}
                      >
                        {r.due_date ? new Date(r.due_date).toLocaleDateString("tr-TR") : "—"}
                        {isOverdue && (
                          <span className="ml-2 text-[10px] bg-red-100 text-red-700 px-1.5 py-0.5 rounded-md uppercase tracking-wider">
                            Gecikti
                          </span>
                        )}
                      </td>

                      <td className="px-5 py-3">
                        {isPending ? (
                          <div className="flex gap-1.5">
                            <button
                              disabled={updating === r.id}
                              onClick={() => updateStatus(r.id, "active")}
                              className="px-2.5 py-1.5 rounded-lg text-xs font-bold cursor-pointer"
                              style={{ backgroundColor: "#D1FAE5", color: "#065F46" }}
                            >
                              Onayla
                            </button>
                            <button
                              disabled={updating === r.id}
                              onClick={() => updateStatus(r.id, "cancelled")}
                              className="px-2.5 py-1.5 rounded-lg text-xs font-bold cursor-pointer"
                              style={{ backgroundColor: "#FEE2E2", color: "#991B1B" }}
                            >
                              Reddet
                            </button>
                          </div>
                        ) : (
                          <select
                            value={r.status}
                            disabled={updating === r.id}
                            onChange={(e) => updateStatus(r.id, e.target.value as Reservation["status"])}
                            className="px-2 py-1.5 rounded-lg text-xs font-semibold outline-none cursor-pointer"
                            style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
                          >
                            <option value="active">Aktif</option>
                            <option value="completed">Tamamlandı</option>
                            <option value="cancelled">İptal Edildi</option>
                            <option value="expired">Süresi Doldu</option>
                          </select>
                        )}
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