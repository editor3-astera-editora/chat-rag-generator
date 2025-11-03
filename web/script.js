// ==============================
// Chat RAG MVP — Script Final
// ==============================

// Envia mensagem para o backend
document.getElementById("chat-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const input = document.getElementById("message-input");
  const message = input.value.trim();
  if (!message) return;

  renderMessage("user", message);
  input.value = "";

  try {
    const res = await fetch("http://localhost:8080/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ user_id: "user123", message }),
    });

    const data = await res.json();
    renderMessage("assistant", data.response, data.sources || []);
  } catch (err) {
    console.error("Erro:", err);
    renderMessage("assistant", "⚠️ Erro ao conectar com o servidor.");
  }
});

// ==============================
// Renderização de mensagens
// ==============================

function renderMessage(role, text, sources = []) {
  const chatBox = document.getElementById("chat-box");

  const container = document.createElement("div");
  container.className = `message-container ${role}`;

  const bubble = document.createElement("div");
  bubble.className = role === "user" ? "user-message" : "assistant-message";
  bubble.innerHTML = renderMarkdownAndMath(text);
  container.appendChild(bubble);

  if (role === "assistant" && Array.isArray(sources) && sources.length > 0) {
    const sourceBlock = document.createElement("div");
    sourceBlock.className = "source-message mt-2";
    sourceBlock.innerHTML = sources.map((s) =>
      `<strong>${s.BookName}</strong> — Unidade ${s.Unit}, Capítulo ${s.Chapter}
       <span class="text-gray-500 ml-1">(${(s.Similarity * 100).toFixed(1)}%)</span>`
    ).join("<br>");
    container.appendChild(sourceBlock);
  }

  chatBox.appendChild(container);
  chatBox.scrollTop = chatBox.scrollHeight;

  // Faz o typeset só do balão recém-adicionado e “apara” <br> ao redor de blocos
  typesetAndTidy(bubble);
}

// ==============================
// Markdown + MathJax Renderer
// ==============================

function renderMarkdownAndMath(text) {
  if (!text) return "";

  text = text.replace(/\(\s*\\\[((?:.|\n){1,60}?)\\\]\s*\)/g, '($\\($1\\)$)');

  text = text.replace(/\\\[((?:.|\n){1,25}?)\\\]/g, '\\($1\\)');

  text = text
    .replace(/\\\(\s+/g, '\\(')
    .replace(/\s+\\\)/g, '\\)')
    .replace(/\\\[\s+/g, '\\[')
    .replace(/\s+\\\]/g, '\\]');


  // Substituições de Markdown básicas
  let html = text
    // títulos
    .replace(/^### (.*$)/gim, '<h3 class="text-lg font-semibold mt-2 mb-1">$1</h3>')
    .replace(/^## (.*$)/gim, '<h2 class="text-xl font-bold mt-3 mb-2">$1</h2>')
    // negrito **texto**
    .replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>")
    // quebra de linha simples
    .replace(/\n/g, "<br>");

  return html;
}

function tidyMathSpacing(containerEl) {
  const blocks = containerEl.querySelectorAll('mjx-container[display="true"]');
  blocks.forEach(block => {
    // remove <br> imediatamente antes
    let prev = block.previousSibling;
    while (prev && prev.nodeType === 1 && prev.tagName === 'BR') {
      const toRemove = prev; prev = prev.previousSibling; toRemove.remove();
    }
    // remove <br> imediatamente depois
    let next = block.nextSibling;
    while (next && next.nodeType === 1 && next.tagName === 'BR') {
      const toRemove = next; next = next.nextSibling; toRemove.remove();
    }
  });
}

function typesetAndTidy(containerEl) {
  if (window.MathJax && window.MathJax.typesetPromise) {
    window.MathJax.typesetPromise([containerEl]).then(() => {
      tidyMathSpacing(containerEl);
    });
  }
}

