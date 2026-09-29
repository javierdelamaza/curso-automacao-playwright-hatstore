// Massa DINÂMICA: cada teste recebe um usuário novo, então nenhum teste depende de outro.
import { randomUUID } from 'node:crypto';

function digitoVerificadorCpf(numeros) {
  const pesoInicial = numeros.length + 1;
  const soma = numeros.reduce((acc, n, i) => acc + n * (pesoInicial - i), 0);
  const resto = soma % 11;
  return resto < 2 ? 0 : 11 - resto;
}

export function gerarCpf() {
  const base = Array.from({ length: 9 }, () => Math.floor(Math.random() * 10));
  const d1 = digitoVerificadorCpf(base);
  const d2 = digitoVerificadorCpf([...base, d1]);
  const n = [...base, d1, d2].join('');
  return `${n.slice(0, 3)}.${n.slice(3, 6)}.${n.slice(6, 9)}-${n.slice(9)}`;
}

export function gerarEmail(prefixo = 'aluno') {
  return `${prefixo}.${randomUUID().slice(0, 8)}@hatstore.test`;
}

function dataNascimentoComIdade(anos) {
  const data = new Date();
  data.setFullYear(data.getFullYear() - anos);
  return data.toISOString().slice(0, 10);
}

export function novoUsuario(sobrescrever = {}) {
  return {
    nome: 'Aluno Playwright',
    email: gerarEmail(),
    cpf: gerarCpf(),
    dataNascimento: dataNascimentoComIdade(25),
    password: 'Senha@123',
    ...sobrescrever,
  };
}

export function usuarioMenorDeIdade() {
  return novoUsuario({ dataNascimento: dataNascimentoComIdade(17) });
}
