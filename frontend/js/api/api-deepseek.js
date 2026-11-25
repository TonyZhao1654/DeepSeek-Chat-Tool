import { request } from "./http.js";

export function getBalance() {
	return request("/api/balance");
}

export function getModels() {
	return request("/api/models");
}

export function openChatStream(body, signal) {
	return fetch("/api/chat", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(body),
		signal,
	});
}
