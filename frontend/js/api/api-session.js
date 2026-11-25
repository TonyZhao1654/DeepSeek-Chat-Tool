import { request } from "./http.js";

export function listSessions() {
	return request("/api/sessions");
}

export function createSession(input) {
	return request("/api/sessions", {
		method: "POST",
		body: JSON.stringify(input),
	});
}

export function getSession(id) {
	return request(`/api/sessions/${encodeURIComponent(String(id))}`);
}

export function updateSession(id, patch) {
	return request(`/api/sessions/${encodeURIComponent(String(id))}`, {
		method: "PUT",
		body: JSON.stringify(patch),
	});
}

export function deleteSession(id) {
	return request(`/api/sessions/${encodeURIComponent(String(id))}`, {
		method: "DELETE",
	});
}
