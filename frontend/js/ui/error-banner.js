import { errorBanner } from "./elements.js";

let errorTimer = null;

export function showError(message) {
	errorBanner.textContent = message;
	errorBanner.hidden = false;
	clearTimeout(errorTimer);
	errorTimer = setTimeout(clearError, 8000);
}

export function clearError() {
	clearTimeout(errorTimer);
	errorBanner.hidden = true;
	errorBanner.textContent = "";
}
