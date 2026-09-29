import Ajv from 'ajv';
import addFormats from 'ajv-formats';

const ajv = new Ajv({ allErrors: true });
addFormats(ajv);

// Retorna a lista de erros legível para aparecer na mensagem do expect.
export function validarContrato(schema, corpo) {
  const validar = ajv.compile(schema);
  const valido = validar(corpo);
  return valido ? [] : validar.errors.map((e) => `${e.instancePath || '(raiz)'} ${e.message}`);
}
