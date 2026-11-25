export async function readChatStream(response, onEvent) {
	const reader = response.body.getReader();
	const decoder = new TextDecoder();
	let buffer = "";

	const emit = (raw) => {
		const parsed = parseEvent(raw);
		if (parsed) onEvent(parsed);
	};

	while (true) {
		const { value, done } = await reader.read();
		if (done) break;
		buffer += decoder.decode(value, { stream: true });
		let idx = buffer.indexOf("\n\n");
		while (idx !== -1) {
			emit(buffer.slice(0, idx));
			buffer = buffer.slice(idx + 2);
			idx = buffer.indexOf("\n\n");
		}
	}
	if (buffer.trim()) emit(buffer);
}

function parseEvent(rawEvent) {
	const lines = rawEvent.split("\n");
	let eventName = "";
	const dataLines = [];
	for (const line of lines) {
		if (line.startsWith("event:")) eventName = line.slice(6).trim();
		else if (line.startsWith("data:")) dataLines.push(line.slice(5).replace(/^\s/, ""));
	}
	if (!dataLines.length) return null;
	const data = dataLines.join("\n");
	if (eventName === "done") {
		try { return { done: true, data: JSON.parse(data) }; } catch { return null; }
	}
	if (data.trim() === "[DONE]") return null;
	try { return { chunk: JSON.parse(data) }; } catch { return null; }
}
