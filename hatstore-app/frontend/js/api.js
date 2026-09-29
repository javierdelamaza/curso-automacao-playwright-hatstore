// Função centralizada para chamadas à API
export async function fetchApi(path, options = {}) {
    const token = localStorage.getItem('jwt_token');
    const userEmail = localStorage.getItem('user_email');

    const headers = {
        'Content-Type': 'application/json',
        'X-Source': 'frontend',
    };

    if (token) headers['Authorization'] = `Bearer ${token}`;
    if (userEmail) headers['X-User-Email'] = userEmail;

    options.headers = { ...headers, ...options.headers };

    const url = path.startsWith('/auth') ? path : `/api${path}`;
    const response = await fetch(url, options);

    if (response.status === 401) {
        localStorage.removeItem('jwt_token');

        const isTest = typeof process !== 'undefined' && process.env && process.env.JEST_WORKER_ID;
        if (!isTest && !window.location.pathname.includes('auth.html')) {
            window.location.href = 'auth.html';
        }
        throw new Error('Não autorizado');
    }

    if (response.status === 503) {
        throw new Error('O site está em manutenção. Tente novamente mais tarde.');
    }

    if (!response.ok) {
        const errorText = (await response.text()).trim();
        throw new Error(errorText || 'Erro na API');
    }
    
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        return response.json();
    }
    return {};
}

// Escapa caracteres HTML para prevenir XSS
export function escapeHtml(str) {
    if (typeof str !== 'string') return str;
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}