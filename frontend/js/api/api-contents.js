import { request } from "./http.js";

export function createContent(sessionId, input) {
	return request(`/api/sessions/${encodeURIComponent(String(sessionId))}/contents`, {
		method: "POST",
		body: JSON.stringify(input),
	});
}

export function updateContent(sessionId, contentId, input) {
	return request(
		`/api/sessions/${encodeURIComponent(String(sessionId))}/contents/${encodeURIComponent(String(contentId))}`,
		{ method: "PUT", body: JSON.stringify(input) }
	);
}

export function deleteContent(sessionId, contentId) {
	return request(
		`/api/sessions/${encodeURIComponent(String(sessionId))}/contents/${encodeURIComponent(String(contentId))}`,
		{ method: "DELETE" }
	);
}
