import type { User } from "../api";

interface SidebarProps {
  user: User;
  active: string;
  onNavigate: (page: string) => void;
  onLogout: () => void;
}

const userNavItems = [
  { key: "books", label: "Kitap Kataloğu", icon: BookIcon },
  { key: "reservations", label: "Rezervasyonlarım", icon: ReservationIcon },
  { key: "chat", label: "Yapay Zeka Asistanı", icon: ChatIcon },
];

const adminNavItems = [
  { key: "books", label: "Kitap Kataloğu", icon: BookIcon },
  { key: "manage-books", label: "Kitap Yönetimi", icon: ManageIcon },
  { key: "manage-users", label: "Kullanıcılar", icon: UsersIcon },
  { key: "all-reservations", label: "Tüm Rezervasyonlar", icon: ReservationIcon },
  { key: "chat", label: "Yapay Zeka Asistanı", icon: ChatIcon },
  { key: "reservations", label: "Rezervasyonlarım", icon: ReservationIcon },
];

export default function Sidebar({ user, active, onNavigate, onLogout }: SidebarProps) {
  const items = user.role === "admin" ? adminNavItems : userNavItems;

  return (
    <aside
      style={{ backgroundColor: "var(--sidebar)", color: "var(--sidebar-foreground)" }}
      className="flex flex-col w-64 shrink-0 h-full"
    >
      <div className="px-6 py-8 border-b border-white/10">
        <div className="flex items-center gap-3">
          <div
            style={{ backgroundColor: "var(--sidebar-accent)", color: "#12110F" }}
            className="w-9 h-9 rounded-lg flex items-center justify-center font-display text-lg font-semibold"
          >
            L
          </div>
          <div>
            <p className="font-display text-lg leading-tight" style={{ color: "var(--sidebar-foreground)" }}>
              Libra
            </p>
            <p className="text-xs opacity-40 leading-tight">Library System</p>
          </div>
        </div>
      </div>

      <nav className="flex-1 px-3 py-5 space-y-0.5 overflow-y-auto">
        {items.map(({ key, label, icon: Icon }) => (
          <button
            key={key}
            onClick={() => onNavigate(key)}
            className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150 cursor-pointer"
            style={{
              backgroundColor: active === key ? "var(--sidebar-muted)" : "transparent",
              color: active === key ? "var(--sidebar-foreground)" : "rgba(232,230,224,0.55)",
            }}
          >
            <Icon size={17} active={active === key} />
            {label}
            {active === key && (
              <span className="ml-auto w-1 h-4 rounded-full" style={{ backgroundColor: "var(--sidebar-accent)" }} />
            )}
          </button>
        ))}
      </nav>

      <div className="px-3 py-4 border-t border-white/10">
        <button
          onClick={() => onNavigate("profile")}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg mb-1 cursor-pointer transition-all duration-150"
          style={{
            backgroundColor: active === "profile" ? "var(--sidebar-accent)" : "var(--sidebar-muted)",
          }}
        >
          <div
            className="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold uppercase shrink-0"
            style={{
              backgroundColor: active === "profile" ? "#12110F" : "var(--sidebar-accent)",
              color: active === "profile" ? "var(--sidebar-accent)" : "#12110F",
            }}
          >
            {user.first_name?.[0] ?? "U"}
          </div>
          <div className="flex-1 min-w-0 text-left">
            <p
              className="text-sm font-semibold truncate"
              style={{ color: active === "profile" ? "#12110F" : "var(--sidebar-foreground)" }}
            >
              {user.first_name}
            </p>
            <p
              className="text-xs capitalize"
              style={{ color: active === "profile" ? "rgba(18,17,15,0.6)" : "rgba(232,230,224,0.4)" }}
            >
              {user.role}
            </p>
          </div>
        </button>
        <button
          onClick={onLogout}
          className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-all duration-150 cursor-pointer"
          style={{ color: "rgba(232,230,224,0.45)" }}
        >
          <LogoutIcon size={16} />
          Sign out
        </button>
      </div>
    </aside>
  );
}

function BookIcon({ size, active }: { size: number; active?: boolean }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={active ? 2.5 : 2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" />
      <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" />
    </svg>
  );
}

function ManageIcon({ size, active }: { size: number; active?: boolean }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={active ? 2.5 : 2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M11 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-6" />
      <path d="M17.5 2.5a2.121 2.121 0 0 1 3 3L12 14l-4 1 1-4z" />
    </svg>
  );
}

function UsersIcon({ size, active }: { size: number; active?: boolean }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={active ? 2.5 : 2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" />
    </svg>
  );
}

function ReservationIcon({ size, active }: { size: number; active?: boolean }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={active ? 2.5 : 2} strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
      <line x1="16" y1="2" x2="16" y2="6" />
      <line x1="8" y1="2" x2="8" y2="6" />
      <line x1="3" y1="10" x2="21" y2="10" />
    </svg>
  );
}

function LogoutIcon({ size }: { size: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
      <polyline points="16 17 21 12 16 7" />
      <line x1="21" y1="12" x2="9" y2="12" />
    </svg>
  );
}

function ChatIcon({ size, active }: { size: number; active?: boolean }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={active ? 2.5 : 2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
    </svg>
  );
}