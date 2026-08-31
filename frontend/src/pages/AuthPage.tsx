import { useState } from "react";
import { api } from "../api";

interface AuthPageProps {
  onLogin: (user: import("../api").User, token: string) => void;
}

export default function AuthPage({ onLogin }: AuthPageProps) {
  const [mode, setMode] = useState<"login" | "register">("login");
  
  const [form, setForm] = useState({ 
    email: "", password: "", first_name: "", last_name: "", phone: "" 
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState("");

  const set = (k: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((f) => ({ ...f, [k]: e.target.value }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    setLoading(true);
    try {
      if (mode === "login") {
        const res = await api.login({ email: form.email, password: form.password });
        
        let tokenString = "";
        if (typeof res === "string") {
          tokenString = res;
        } else if (res && typeof res === "object") {
          tokenString = res.token || res.data || Object.values(res)[0] as string; 
        }

        if (!tokenString || typeof tokenString !== "string") {
          throw new Error("Token bulunamadı, backend yanıtını kontrol edin.");
        }

        let base64 = tokenString.split('.')[0].replace(/-/g, '+').replace(/_/g, '/');
        while (base64.length % 4) base64 += '=';
        const payload = JSON.parse(window.atob(base64));

        const loggedInUser = {
          id: payload.user_id,
          role: payload.role,
          email: form.email,
        };

        onLogin(loggedInUser as import("../api").User, tokenString);
        
      } else {
        await api.register({ 
          email: form.email, 
          password: form.password,
          first_name: form.first_name,
          last_name: form.last_name,
          phone: form.phone
        });
        setSuccess("Kayıt başarılı! Şimdi giriş yapabilirsiniz.");
        setMode("login");
        setForm({ email: "", password: "", first_name: "", last_name: "", phone: "" });
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Bir hata oluştu");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex" style={{ backgroundColor: "var(--background)" }}>
      {/* Sol Panel */}
      <div className="hidden lg:flex lg:w-[45%] flex-col justify-between p-14 relative overflow-hidden" style={{ backgroundColor: "var(--sidebar)" }}>
        <div className="absolute inset-0 pointer-events-none select-none" aria-hidden>
          {[{ x: "10%", y: "18%", w: 38, h: 54, rot: -18, opacity: 0.06 }, { x: "22%", y: "28%", w: 28, h: 42, rot: 10, opacity: 0.05 }, { x: "68%", y: "10%", w: 44, h: 62, rot: 14, opacity: 0.07 }, { x: "80%", y: "30%", w: 30, h: 44, rot: -8, opacity: 0.05 }, { x: "55%", y: "70%", w: 36, h: 52, rot: 20, opacity: 0.06 }, { x: "15%", y: "72%", w: 24, h: 36, rot: -12, opacity: 0.04 }, { x: "40%", y: "82%", w: 32, h: 46, rot: 6, opacity: 0.05 }, { x: "75%", y: "75%", w: 22, h: 34, rot: -20, opacity: 0.04 }].map((b, i) => (
            <div key={i} className="absolute rounded-sm" style={{ left: b.x, top: b.y, width: b.w, height: b.h, transform: `rotate(${b.rot}deg)`, backgroundColor: "var(--sidebar-foreground)", opacity: b.opacity }} />
          ))}
        </div>
        <div className="relative flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl flex items-center justify-center font-display text-xl" style={{ backgroundColor: "var(--sidebar-accent)", color: "#12110F" }}>L</div>
          <span className="font-display text-2xl" style={{ color: "var(--sidebar-foreground)" }}>Libra</span>
        </div>
        <div className="relative">
          <div className="w-12 h-px mb-8" style={{ backgroundColor: "var(--sidebar-accent)" }} />
          <h2 className="font-display text-5xl leading-[1.1] mb-6" style={{ color: "var(--sidebar-foreground)" }}>Okumak<br />bir <span style={{ color: "var(--sidebar-accent)", fontStyle: "italic" }}>ayrıcalık</span><br />değil, haktır.</h2>
          <p className="text-sm leading-relaxed max-w-xs" style={{ color: "rgba(232,230,224,0.45)" }}>Dijital kütüphane sisteminizle binlerce esere ulaşın, rezervasyon yapın, okuma geçmişinizi takip edin.</p>
        </div>
        <div className="relative flex items-center gap-6">
          {[{ label: "kitap koleksiyonu", value: "12,400+" }, { label: "aktif üye", value: "3,200+" }].map(({ label, value }) => (
            <div key={label}>
              <p className="font-display text-3xl mb-0.5" style={{ color: "var(--sidebar-foreground)" }}>{value}</p>
              <p className="text-xs uppercase tracking-widest" style={{ color: "rgba(232,230,224,0.35)" }}>{label}</p>
            </div>
          ))}
          <div className="w-px h-10 mx-2" style={{ backgroundColor: "rgba(232,230,224,0.12)" }} />
          <div>
            <p className="font-display text-3xl mb-0.5" style={{ color: "var(--sidebar-foreground)" }}>48K+</p>
            <p className="text-xs uppercase tracking-widest" style={{ color: "rgba(232,230,224,0.35)" }}>rezervasyon</p>
          </div>
        </div>
      </div>

      {/* Sağ Panel */}
      <div className="flex-1 flex items-center justify-center px-6 py-12 overflow-y-auto">
        <div className="w-full max-w-md">
          <div className="mb-8">
            <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
              {mode === "login" ? "Tekrar hoş geldiniz" : "Hesap oluşturun"}
            </h1>
            <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
              {mode === "login" ? "Hesabınıza giriş yapın" : "Yeni bir hesap oluşturun"}
            </p>
          </div>

          <div className="flex rounded-xl p-1 mb-8" style={{ backgroundColor: "var(--muted)" }}>
            {(["login", "register"] as const).map((m) => (
              <button
                key={m}
                onClick={() => { setMode(m); setError(""); setSuccess(""); }}
                className="flex-1 py-2 rounded-lg text-sm font-semibold transition-all duration-200 cursor-pointer"
                style={{ backgroundColor: mode === m ? "var(--card)" : "transparent", color: mode === m ? "var(--foreground)" : "var(--muted-foreground)", boxShadow: mode === m ? "0 1px 3px rgba(0,0,0,0.08)" : "none" }}
              >
                {m === "login" ? "Giriş Yap" : "Kayıt Ol"}
              </button>
            ))}
          </div>

          {error && <div className="mb-4 px-4 py-3 rounded-lg text-sm font-medium" style={{ backgroundColor: "#FEF2F2", color: "#DC2626", border: "1px solid #FECACA" }}>{error}</div>}
          {success && <div className="mb-4 px-4 py-3 rounded-lg text-sm font-medium" style={{ backgroundColor: "#F0FDF4", color: "#16A34A", border: "1px solid #BBF7D0" }}>{success}</div>}

          <form onSubmit={handleSubmit} className="space-y-4">
            {mode === "register" && (
              <>
                <div className="flex gap-4">
                  <div className="flex-1">
                    <Field label="Ad" type="text" value={form.first_name} onChange={set("first_name")} placeholder="Furkan" required />
                  </div>
                  <div className="flex-1">
                    <Field label="Soyad" type="text" value={form.last_name} onChange={set("last_name")} placeholder="Yılmaz" required />
                  </div>
                </div>
                {/* Kullanıcı Adı alanı buradan silindi */}
                <Field label="Telefon" type="tel" value={form.phone} onChange={set("phone")} placeholder="0555..." required />
              </>
            )}
            <Field label="E-posta" type="email" value={form.email} onChange={set("email")} placeholder="ornek@email.com" required />
            <Field label="Şifre" type="password" value={form.password} onChange={set("password")} placeholder="••••••••" required />

            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 rounded-xl text-sm font-bold tracking-wide transition-all duration-200 cursor-pointer mt-2"
              style={{ backgroundColor: loading ? "var(--muted)" : "var(--primary)", color: loading ? "var(--muted-foreground)" : "var(--primary-foreground)" }}
            >
              {loading ? "Lütfen bekleyin…" : mode === "login" ? "Giriş Yap" : "Hesap Oluştur"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

function Field({ label, type, value, onChange, placeholder, required }: { label: string; type: string; value: string; onChange: (e: React.ChangeEvent<HTMLInputElement>) => void; placeholder: string; required?: boolean; }) {
  return (
    <div>
      <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>{label}</label>
      <input type={type} value={value} onChange={onChange} placeholder={placeholder} required={required} className="w-full px-4 py-2.5 rounded-xl text-sm outline-none transition-all duration-150" style={{ backgroundColor: "var(--card)", border: "1.5px solid var(--border)", color: "var(--foreground)" }} onFocus={(e) => (e.currentTarget.style.borderColor = "var(--primary)")} onBlur={(e) => (e.currentTarget.style.borderColor = "var(--border)")} />
    </div>
  );
}