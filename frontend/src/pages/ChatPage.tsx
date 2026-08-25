import { useState, useRef, useEffect } from "react";

interface Message {
  id: number;
  role: "user" | "assistant";
  content: string;
  ts: Date;
}

const SUGGESTIONS = [
  "Hangi kitaplar mevcut?",
  "Nasıl rezervasyon yapabilirim?",
  "Popüler roman önerir misin?",
  "Rezervasyonumu nasıl iptal ederim?",
];

export default function ChatPage() {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: 0,
      role: "assistant",
      content: "Merhaba! Kütüphane asistanınım. Kitap önerileri, rezervasyon bilgileri veya koleksiyon hakkında sormak istediğiniz her şeyi sorabilirsiniz.",
      ts: new Date(),
    },
  ]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLTextAreaElement>(null);
  const idRef = useRef(1);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, loading]);

  const send = async (text: string) => {
    const content = text.trim();
    if (!content || loading) return;
    setInput("");

    const userMsg: Message = { id: idRef.current++, role: "user", content, ts: new Date() };
    setMessages((ms) => [...ms, userMsg]);
    setLoading(true);

    try {
      const res = await fetch("/api/chat", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
        },
        body: JSON.stringify({ message: content }),
      });
      if (!res.ok) throw new Error(await res.text());
      const data = await res.json();
      const reply = data.reply ?? data.message ?? data.response ?? data[""] ?? Object.values(data)[0] ?? JSON.stringify(data);
      setMessages((ms) => [...ms, { id: idRef.current++, role: "assistant", content: reply, ts: new Date() }]);
    } catch (e: unknown) {
      setMessages((ms) => [
        ...ms,
        {
          id: idRef.current++,
          role: "assistant",
          content: e instanceof Error ? `Hata: ${e.message}` : "Bir hata oluştu, tekrar deneyin.",
          ts: new Date(),
        },
      ]);
    } finally {
      setLoading(false);
      inputRef.current?.focus();
    }
  };

  const handleKey = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      send(input);
    }
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div
        className="px-8 py-5 flex items-center gap-4 shrink-0"
        style={{ borderBottom: "1px solid var(--border)", backgroundColor: "var(--card)" }}
      >
        <div
          className="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
          style={{ backgroundColor: "var(--sidebar)", color: "var(--sidebar-accent)" }}
        >
          <BotIcon />
        </div>
        <div>
          <h1 className="font-display text-xl leading-tight" style={{ color: "var(--foreground)" }}>
            Kütüphane Asistanı
          </h1>
          <div className="flex items-center gap-1.5 mt-0.5">
            <span className="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse" />
            <span className="text-xs" style={{ color: "var(--muted-foreground)" }}>Çevrimiçi</span>
          </div>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto px-6 py-6 space-y-5">
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={`flex gap-3 ${msg.role === "user" ? "flex-row-reverse" : "flex-row"}`}
          >
            {msg.role === "assistant" && (
              <div
                className="w-8 h-8 rounded-xl flex items-center justify-center shrink-0 mt-0.5"
                style={{ backgroundColor: "var(--sidebar)", color: "var(--sidebar-accent)" }}
              >
                <BotIcon size={16} />
              </div>
            )}
            <div className={`max-w-[72%] flex flex-col gap-1 ${msg.role === "user" ? "items-end" : "items-start"}`}>
              <div
                className="px-4 py-3 rounded-2xl text-sm leading-relaxed whitespace-pre-wrap"
                style={
                  msg.role === "user"
                    ? { backgroundColor: "var(--primary)", color: "var(--primary-foreground)", borderBottomRightRadius: 6 }
                    : { backgroundColor: "var(--card)", color: "var(--foreground)", border: "1px solid var(--border)", borderBottomLeftRadius: 6 }
                }
              >
                {msg.content}
              </div>
              <span className="text-xs px-1" style={{ color: "var(--muted-foreground)" }}>
                {msg.ts.toLocaleTimeString("tr-TR", { hour: "2-digit", minute: "2-digit" })}
              </span>
            </div>
          </div>
        ))}

        {loading && (
          <div className="flex gap-3">
            <div
              className="w-8 h-8 rounded-xl flex items-center justify-center shrink-0 mt-0.5"
              style={{ backgroundColor: "var(--sidebar)", color: "var(--sidebar-accent)" }}
            >
              <BotIcon size={16} />
            </div>
            <div
              className="px-4 py-3.5 rounded-2xl flex items-center gap-1.5"
              style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)", borderBottomLeftRadius: 6 }}
            >
              {[0, 1, 2].map((i) => (
                <span
                  key={i}
                  className="w-1.5 h-1.5 rounded-full"
                  style={{
                    backgroundColor: "var(--muted-foreground)",
                    animation: `bounce 1.2s ease-in-out ${i * 0.2}s infinite`,
                  }}
                />
              ))}
            </div>
          </div>
        )}

        <div ref={bottomRef} />
      </div>

      {/* Suggestions — shown only at start */}
      {messages.length === 1 && (
        <div className="px-6 pb-3 flex flex-wrap gap-2">
          {SUGGESTIONS.map((s) => (
            <button
              key={s}
              onClick={() => send(s)}
              className="px-3 py-2 rounded-xl text-xs font-semibold cursor-pointer transition-all duration-150 hover:-translate-y-0.5"
              style={{
                backgroundColor: "var(--card)",
                border: "1.5px solid var(--border)",
                color: "var(--foreground)",
              }}
            >
              {s}
            </button>
          ))}
        </div>
      )}

      {/* Input */}
      <div
        className="px-6 py-4 shrink-0"
        style={{ borderTop: "1px solid var(--border)", backgroundColor: "var(--card)" }}
      >
        <div
          className="flex items-end gap-3 rounded-2xl px-4 py-3"
          style={{ backgroundColor: "var(--background)", border: "1.5px solid var(--border)" }}
        >
          <textarea
            ref={inputRef}
            rows={1}
            value={input}
            onChange={(e) => {
              setInput(e.target.value);
              e.target.style.height = "auto";
              e.target.style.height = Math.min(e.target.scrollHeight, 120) + "px";
            }}
            onKeyDown={handleKey}
            placeholder="Bir şey sorun… (Enter ile gönderin)"
            className="flex-1 resize-none outline-none text-sm leading-relaxed bg-transparent"
            style={{ color: "var(--foreground)", maxHeight: 120 }}
          />
          <button
            onClick={() => send(input)}
            disabled={!input.trim() || loading}
            className="w-9 h-9 rounded-xl flex items-center justify-center shrink-0 transition-all duration-150 cursor-pointer"
            style={{
              backgroundColor: input.trim() && !loading ? "var(--primary)" : "var(--muted)",
              color: input.trim() && !loading ? "var(--primary-foreground)" : "var(--muted-foreground)",
            }}
          >
            <SendIcon />
          </button>
        </div>
        <p className="text-xs text-center mt-2" style={{ color: "var(--muted-foreground)", opacity: 0.6 }}>
          Shift+Enter ile yeni satır ekleyin
        </p>
      </div>

      <style>{`
        @keyframes bounce {
          0%, 60%, 100% { transform: translateY(0); }
          30% { transform: translateY(-5px); }
        }
      `}</style>
    </div>
  );
}

function BotIcon({ size = 18 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="11" width="18" height="10" rx="2" />
      <circle cx="12" cy="5" r="2" />
      <path d="M12 7v4" />
      <line x1="8" y1="16" x2="8" y2="16" strokeWidth="2.5" />
      <line x1="12" y1="16" x2="12" y2="16" strokeWidth="2.5" />
      <line x1="16" y1="16" x2="16" y2="16" strokeWidth="2.5" />
    </svg>
  );
}

function SendIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
      <line x1="22" y1="2" x2="11" y2="13" />
      <polygon points="22 2 15 22 11 13 2 9 22 2" />
    </svg>
  );
}
