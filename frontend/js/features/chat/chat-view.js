import { chatElement } from "../../ui/elements.js";

export function renderMessages(messages) {
	chatElement.innerHTML = "";
	(messages || []).forEach((m) => {
		if (!m.content) return;
		if (m.role === "user") appendMessage("user", m.content);
		else if (m.role === "assistant") appendMessage("ai", m.content, true);
	});
	scrollToBottom();
}

export function createAssistantBubble() {
	const wrapper = document.createElement("div");
	wrapper.className = "message ai";

	const reasoning = document.createElement("div");
	reasoning.className = "reasoning";
	reasoning.hidden = true;

	const thinking = document.createElement("div");
	thinking.className = "thinking";
	thinking.innerHTML = `
		<div class="thinking-dot"></div>
		<div class="thinking-dot"></div>
		<div class="thinking-dot"></div>
		<span>正在思考...</span>
	`;

	const content = document.createElement("div");
	content.className = "content";

	wrapper.append(reasoning, thinking, content);
	chatElement.appendChild(wrapper);
	scrollToBottom();
	return { wrapper, thinking, reasoning, content };
}

export function renderAssistantContent(bubble, text) {
	if (bubble.thinking.isConnected) bubble.thinking.remove();
	bubble.content.innerHTML = renderMarkdown(text);
	scrollToBottom();
}

export function appendMessage(role, text, isMarkdown = false) {
	const el = document.createElement("div");
	el.className = `message ${role}`;
	if (role === "meta") {
		el.classList.add("meta");
		el.textContent = text;
	} else if (isMarkdown) {
		el.innerHTML = renderMarkdown(text);
	} else {
		el.textContent = text;
	}
	chatElement.appendChild(el);
	scrollToBottom();
	return el;
}

export function appendMeta(text) {
	appendMessage("meta", text);
}

export function scrollToBottom() {
	chatElement.scrollTop = chatElement.scrollHeight;
}

function renderMarkdown(content) {
	if (window.marked) return window.marked.parse(content);
	const el = document.createElement("div");
	el.textContent = content;
	return el.innerHTML;
}
