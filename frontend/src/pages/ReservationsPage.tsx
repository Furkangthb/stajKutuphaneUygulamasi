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

interface ReservationsPageProps {
  userId: number;
}

export default function ReservationsPage({ userId }: ReservationsPageProps) {
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [updating, setUpdating] = useState<number | null>(null);
  const [toast, setToast] = useState("");

  useEffect(() => {
    api.getUserReservations(userId)
      .then(setReservations)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [userId]);

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

  return (
    <div className="h-full flex flex-col">
      <div className="px-8 pt-8 pb-6">
        <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
          Rezervasyonlarım
        </h1>
        <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
          {reservations.length} rezervasyon
        </p>
      </div>

      <div className="flex-1 overflow-y-auto px-8 pb-8">
        {loading ? (
          <LoadingList />
        ) : error ? (
          <div className="flex items-center justify-center h-48">
            <p className="text-sm" style={{ color: "#DC2626" }}>{error}</p>
          </div>
        ) : reservations.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-48 gap-2">
            <p className="text-sm font-semibold" style={{ color: "var(--foreground)" }}>Henüz rezervasyonunuz yok</p>
            <p className="text-xs" style={{ color: "var(--muted-foreground)" }}>Kitap kataloğundan rezervasyon yapabilirsiniz</p>
          </div>
        ) : (
          <div className="space-y-3 max-w-2xl">
            {reservations.map((res) => (
              <ReservationCard
                key={res.id}
                reservation={res}
                onUpdate={updateStatus}
                updating={updating === res.id}
              />
            ))}
          </div>
        )}
      </div>

      {toast && (
        <div
          className="fixed bottom-6 right-6 px-5 py-3 rounded-xl text-sm font-semibold shadow-lg z-50"
          style={{ backgroundColor: "var(--sidebar)", color: "var(--sidebar-foreground)" }}
        >
          {toast}
        </div>
      )}
    </div>
  );
}

function ReservationCard({
  reservation, onUpdate, updating,
}: {
  reservation: Reservation;
  onUpdate: (id: number, status: Reservation["status"]) => void;
  updating: boolean;
}) {
  const sc = STATUS_COLORS[reservation.status] ?? { bg: "#F3F4F6", text: "#6B7280" };
  const title = reservation.book?.title ?? `Kitap #${reservation.book_id}`;
  const author = reservation.book?.author ?? "";

  return (
    <div
      className="rounded-2xl p-5 flex items-start gap-4"
      style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}
    >
      <div
        className="w-12 h-16 rounded-lg flex items-center justify-center shrink-0"
        style={{ backgroundColor: "var(--muted)" }}
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="var(--muted-foreground)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
          <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" />
          <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" />
        </svg>
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2 mb-1">
          <h3 className="font-semibold text-sm leading-tight" style={{ color: "var(--foreground)" }}>{title}</h3>
          <span
            className="shrink-0 px-2 py-0.5 rounded-full text-xs font-bold"
            style={{ backgroundColor: sc.bg, color: sc.text }}
          >
            {STATUS_LABELS[reservation.status]}
          </span>
        </div>
        {author && <p className="text-xs mb-2" style={{ color: "var(--muted-foreground)" }}>{author}</p>}
        {reservation.reserved_at && (
          <p className="text-xs" style={{ color: "var(--muted-foreground)" }}>
            Rezervasyon: {new Date(reservation.reserved_at).toLocaleDateString("tr-TR")}
          </p>
        )}
        {reservation.status === "active" && (
          <div className="flex gap-2 mt-3">
            <button
              disabled={updating}
              onClick={() => onUpdate(reservation.id, "returned")}
              className="px-3 py-1.5 rounded-lg text-xs font-semibold cursor-pointer transition-all"
              style={{ backgroundColor: "var(--secondary)", color: "var(--secondary-foreground)" }}
            >
              {updating ? "…" : "İade Et"}
            </button>
            <button
              disabled={updating}
              onClick={() => onUpdate(reservation.id, "cancelled")}
              className="px-3 py-1.5 rounded-lg text-xs font-semibold cursor-pointer transition-all"
              style={{ backgroundColor: "#FEE2E2", color: "#991B1B" }}
            >
              {updating ? "…" : "İptal Et"}
            </button>
          </div>
        )}
        {reservation.status === "pending" && (
          <button
            disabled={updating}
            onClick={() => onUpdate(reservation.id, "cancelled")}
            className="mt-3 px-3 py-1.5 rounded-lg text-xs font-semibold cursor-pointer transition-all"
            style={{ backgroundColor: "#FEE2E2", color: "#991B1B" }}
          >
            {updating ? "…" : "İptal Et"}
          </button>
        )}
      </div>
    </div>
  );
}

function LoadingList() {
  return (
    <div className="space-y-3 max-w-2xl">
      {Array.from({ length: 4 }).map((_, i) => (
        <div
          key={i}
          className="rounded-2xl p-5 flex gap-4 animate-pulse"
          style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}
        >
          <div className="w-12 h-16 rounded-lg" style={{ backgroundColor: "var(--muted)" }} />
          <div className="flex-1 space-y-2">
            <div className="h-4 rounded" style={{ backgroundColor: "var(--muted)", width: "60%" }} />
            <div className="h-3 rounded" style={{ backgroundColor: "var(--muted)", width: "40%" }} />
          </div>
        </div>
      ))}
    </div>
  );
}
