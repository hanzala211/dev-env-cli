import "./App.css";
import { useEffect, useState } from "react";

type Project = {
  name: string;
  cmd: string;
  path: string;
};

type UIProject = Project & { running: boolean };

const pageStyle: React.CSSProperties = {
  padding: 24,
  maxWidth: 980,
  margin: "0 auto",
  fontFamily:
    "ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial",
};

const cardStyle: React.CSSProperties = {
  backgroundColor: "#ffffff",
  border: "1px solid #e5e7eb",
  borderRadius: 10,
  boxShadow: "0 1px 2px rgba(0,0,0,0.04)",
  padding: 16,
};

const tableStyle: React.CSSProperties = {
  width: "100%",
  borderCollapse: "collapse",
  fontFamily:
    "ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial",
  marginTop: 16,
};

const cellStyle: React.CSSProperties = {
  border: "1px solid #e5e7eb",
  padding: "10px 12px",
  textAlign: "left",
  verticalAlign: "top",
  fontSize: 14,
  color: "#111827",
};

const headerCellStyle: React.CSSProperties = {
  ...cellStyle,
  backgroundColor: "#f9fafb",
  fontWeight: 600,
};

const toolbarStyle: React.CSSProperties = {
  display: "flex",
  gap: 8,
  alignItems: "center",
  justifyContent: "space-between",
  marginTop: 8,
  marginBottom: 8,
};

const buttonBaseStyle: React.CSSProperties = {
  border: "1px solid transparent",
  borderRadius: 6,
  padding: "8px 12px",
  fontSize: 13,
  fontWeight: 600,
  cursor: "pointer",
  transition: "background-color 120ms, border-color 120ms, opacity 120ms",
};

const startBtnStyle: React.CSSProperties = {
  ...buttonBaseStyle,
  backgroundColor: "#065f46",
  borderColor: "#065f46",
  color: "#ffffff",
};

const ghostBtnStyle: React.CSSProperties = {
  ...buttonBaseStyle,
  backgroundColor: "#f3f4f6",
  borderColor: "#e5e7eb",
  color: "#111827",
};

const disabledStyle: React.CSSProperties = {
  opacity: 0.5,
  cursor: "not-allowed",
};

const statusPillStyle = (running: boolean): React.CSSProperties => ({
  display: "inline-block",
  padding: "2px 8px",
  borderRadius: 999,
  fontSize: 12,
  fontWeight: 700,
  color: running ? "#065f46" : "#92400e",
  backgroundColor: running ? "#d1fae5" : "#fef3c7",
  border: `1px solid ${running ? "#34d399" : "#fbbf24"}`,
});

function App() {
  const [projects, setProjects] = useState<UIProject[]>([]);

  useEffect(() => {
    let isMounted = true;
    async function loadProjects() {
      try {
        const res = await fetch("/api/projects");
        if (!res.ok) return;
        const data: { projects: UIProject[] } = await res.json();
        if (isMounted) setProjects(data?.projects ?? []);
      } catch {
        // Swallow network errors silently for now
      }
    }
    loadProjects();
    return () => {
      isMounted = false;
    };
  }, []);

  const handleStart = async (index: number) => {
    const project = projects[index];
    try {
      const res = await fetch("/api/projects/start", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: project.name }),
      });
      if (!res.ok) return;
      setProjects((prev) =>
        prev.map((p, i) => (i === index ? { ...p, running: true } : p))
      );
    } catch (error) {
      console.error(error);
    }
  };

  const handleStop = async (index: number) => {
    const project = projects[index];
    try {
      const res = await fetch("/api/projects/stop", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: project.name }),
      });
      if (!res.ok) return;
      setProjects((prev) =>
        prev.map((p, i) => (i === index ? { ...p, running: false } : p))
      );
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div style={pageStyle}>
      <div
        style={{
          display: "flex",
          alignItems: "baseline",
          justifyContent: "space-between",
        }}
      >
        <h1
          style={{ margin: 0, fontSize: 20, fontWeight: 800, color: "#111827" }}
        >
          Projects
        </h1>
      </div>

      <div style={toolbarStyle}>
        <div style={{ fontSize: 12, color: "#6b7280" }}>
          {projects.filter((p) => p.running).length} running / {projects.length}{" "}
          total
        </div>
      </div>

      <div style={cardStyle}>
        <table style={tableStyle}>
          <thead>
            <tr>
              <th style={headerCellStyle}>Name</th>
              <th style={headerCellStyle}>Command</th>
              <th style={headerCellStyle}>Path</th>
              <th style={headerCellStyle}>Status</th>
              <th style={headerCellStyle}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {projects.map((p, idx) => (
              <tr
                key={p.name}
                style={{
                  backgroundColor: idx % 2 === 0 ? "#ffffff" : "#fcfcfd",
                }}
              >
                <td style={cellStyle}>{p.name}</td>
                <td style={cellStyle}>
                  <code
                    style={{
                      backgroundColor: "#f3f4f6",
                      padding: "2px 6px",
                      borderRadius: 4,
                    }}
                  >
                    {p.cmd}
                  </code>
                </td>
                <td style={cellStyle}>{p.path}</td>
                <td style={cellStyle}>
                  <span style={statusPillStyle(p.running)}>
                    {p.running ? "Running" : "Stopped"}
                  </span>
                </td>
                <td style={{ ...cellStyle, whiteSpace: "nowrap" }}>
                  <button
                    onClick={() => handleStart(idx)}
                    style={{
                      ...startBtnStyle,
                      ...(p.running ? disabledStyle : {}),
                      marginRight: 8,
                    }}
                    disabled={p.running}
                  >
                    Start
                  </button>
                  <button
                    onClick={() => handleStop(idx)}
                    style={{
                      ...ghostBtnStyle,
                      ...(p.running ? {} : disabledStyle),
                      borderColor: p.running ? "#e5e7eb" : undefined,
                    }}
                    disabled={!p.running}
                  >
                    Stop
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default App;
