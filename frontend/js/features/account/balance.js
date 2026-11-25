import { getBalance } from "../../api/api-deepseek.js";
import { balanceValue } from "../../ui/elements.js";
import { showError } from "../../ui/error-banner.js";

export function formatBalance(value) {
	const n = parseFloat(value);
	return Number.isFinite(n) ? n.toFixed(2) : String(value);
}

export function updateBalance(value) {
	balanceValue.textContent = formatBalance(value);
}

export async function refreshBalance() {
	try {
		const balance = await getBalance();
		updateBalance(balance && balance.total_balance);
	} catch (error) {
		balanceValue.textContent = "获取失败";
		showError("余额查询失败：" + error.message);
	}
}
