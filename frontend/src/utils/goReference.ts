import type { Language } from '../i18n';

export interface TypeInfo {
  title: string;
  description: string;
}

export interface FuncInfo {
  signature: string;
  description: string;
  returns?: string;
}

const typesEn: Record<string, TypeInfo> = {
  bool: {
    title: 'bool',
    description: 'Boolean: true or false.',
  },
  string: {
    title: 'string',
    description: 'Immutable UTF-8 text sequence.',
  },
  byte: {
    title: 'byte',
    description: 'Alias for uint8 — one byte of data.',
  },
  rune: {
    title: 'rune',
    description: 'Alias for int32 — a Unicode code point.',
  },
  int: {
    title: 'int',
    description: 'Signed integer (32 or 64 bits depending on platform).',
  },
  int8: { title: 'int8', description: 'Signed 8-bit integer (−128 … 127).' },
  int16: { title: 'int16', description: 'Signed 16-bit integer.' },
  int32: { title: 'int32', description: 'Signed 32-bit integer.' },
  int64: { title: 'int64', description: 'Signed 64-bit integer.' },
  uint: { title: 'uint', description: 'Unsigned integer (platform size).' },
  uint8: { title: 'uint8', description: 'Unsigned 8-bit integer (0 … 255).' },
  uint16: { title: 'uint16', description: 'Unsigned 16-bit integer.' },
  uint32: { title: 'uint32', description: 'Unsigned 32-bit integer.' },
  uint64: { title: 'uint64', description: 'Unsigned 64-bit integer.' },
  uintptr: {
    title: 'uintptr',
    description: 'Unsigned integer large enough to store a pointer.',
  },
  float32: { title: 'float32', description: 'IEEE-754 32-bit floating point.' },
  float64: { title: 'float64', description: 'IEEE-754 64-bit floating point.' },
  complex64: { title: 'complex64', description: 'Complex number with float32 parts.' },
  complex128: { title: 'complex128', description: 'Complex number with float64 parts.' },
  error: {
    title: 'error',
    description: 'Built-in interface for errors: a value with an Error() string method.',
  },
  any: {
    title: 'any',
    description: 'Alias for interface{} — can hold any value.',
  },
};

const typesRu: Record<string, TypeInfo> = {
  bool: { title: 'bool', description: 'Логический тип: true или false.' },
  string: { title: 'string', description: 'Неизменяемая UTF-8 строка.' },
  byte: { title: 'byte', description: 'Псевдоним uint8 — один байт.' },
  rune: { title: 'rune', description: 'Псевдоним int32 — символ Unicode.' },
  int: { title: 'int', description: 'Знаковое целое (32/64 бита — зависит от платформы).' },
  int8: { title: 'int8', description: 'Знаковое 8-битное целое (−128 … 127).' },
  int16: { title: 'int16', description: 'Знаковое 16-битное целое.' },
  int32: { title: 'int32', description: 'Знаковое 32-битное целое.' },
  int64: { title: 'int64', description: 'Знаковое 64-битное целое.' },
  uint: { title: 'uint', description: 'Беззнаковое целое (размер платформы).' },
  uint8: { title: 'uint8', description: 'Беззнаковое 8-битное (0 … 255).' },
  uint16: { title: 'uint16', description: 'Беззнаковое 16-битное целое.' },
  uint32: { title: 'uint32', description: 'Беззнаковое 32-битное целое.' },
  uint64: { title: 'uint64', description: 'Беззнаковое 64-битное целое.' },
  uintptr: { title: 'uintptr', description: 'Целое, достаточное для хранения указателя.' },
  float32: { title: 'float32', description: 'Число с плавающей точкой 32 бита.' },
  float64: { title: 'float64', description: 'Число с плавающей точкой 64 бита.' },
  complex64: { title: 'complex64', description: 'Комплексное число (float32).' },
  complex128: { title: 'complex128', description: 'Комплексное число (float64).' },
  error: {
    title: 'error',
    description: 'Встроенный интерфейс ошибок: значение с методом Error() string.',
  },
  any: { title: 'any', description: 'Псевдоним interface{} — любое значение.' },
};

const stdlibEn: Record<string, FuncInfo> = {
  'fmt.Println': {
    signature: 'func Println(a ...any) (n int, err error)',
    description: 'Prints arguments with spaces, then newline. Returns bytes written.',
  },
  'fmt.Printf': {
    signature: 'func Printf(format string, a ...any) (n int, err error)',
    description: 'Formatted print to stdout.',
  },
  'fmt.Print': {
    signature: 'func Print(a ...any) (n int, err error)',
    description: 'Prints arguments to stdout without trailing newline.',
  },
  'fmt.Sprintf': {
    signature: 'func Sprintf(format string, a ...any) string',
    description: 'Returns formatted string without printing.',
    returns: 'string',
  },
  'len': {
    signature: 'func len(v Type) int',
    description: 'Returns length of array, slice, string, map, or channel buffer.',
    returns: 'int',
  },
  'cap': {
    signature: 'func cap(v Type) int',
    description: 'Capacity of slice or channel buffer.',
    returns: 'int',
  },
  'make': {
    signature: 'func make(T, size ...int) T',
    description: 'Allocates and initializes slice, map, or channel.',
  },
  'new': {
    signature: 'func new(Type) *Type',
    description: 'Allocates zeroed storage and returns pointer to it.',
  },
  'append': {
    signature: 'func append(slice []T, elems ...T) []T',
    description: 'Appends elements to a slice, growing if needed.',
    returns: '[]T',
  },
  'copy': {
    signature: 'func copy(dst, src []T) int',
    description: 'Copies min(len(dst), len(src)) elements. Returns count copied.',
    returns: 'int',
  },
  'delete': {
    signature: 'func delete(m map[K]V, key K)',
    description: 'Removes key from map (no-op if key absent).',
  },
  'panic': {
    signature: 'func panic(v any)',
    description: 'Stops normal execution; starts panicking.',
  },
  'recover': {
    signature: 'func recover() any',
    description: 'Stops panic in deferred call; returns panic value.',
    returns: 'any',
  },
};

const stdlibRu: Record<string, FuncInfo> = {
  'fmt.Println': {
    signature: 'func Println(a ...any) (n int, err error)',
    description: 'Печатает аргументы через пробел и перевод строки.',
  },
  'fmt.Printf': {
    signature: 'func Printf(format string, a ...any) (n int, err error)',
    description: 'Форматированный вывод в stdout.',
  },
  'fmt.Print': {
    signature: 'func Print(a ...any) (n int, err error)',
    description: 'Печать без перевода строки в конце.',
  },
  'fmt.Sprintf': {
    signature: 'func Sprintf(format string, a ...any) string',
    description: 'Возвращает отформатированную строку.',
    returns: 'string',
  },
  len: {
    signature: 'func len(v Type) int',
    description: 'Длина массива, слайса, строки, map или буфера канала.',
    returns: 'int',
  },
  cap: {
    signature: 'func cap(v Type) int',
    description: 'Ёмкость слайса или буфера канала.',
    returns: 'int',
  },
  make: {
    signature: 'func make(T, size ...int) T',
    description: 'Создаёт и инициализирует slice, map или channel.',
  },
  new: {
    signature: 'func new(Type) *Type',
    description: 'Выделяет память под нулевое значение и возвращает указатель.',
  },
  append: {
    signature: 'func append(slice []T, elems ...T) []T',
    description: 'Добавляет элементы в слайс, при необходимости расширяя его.',
    returns: '[]T',
  },
  copy: {
    signature: 'func copy(dst, src []T) int',
    description: 'Копирует min(len(dst), len(src)) элементов.',
    returns: 'int',
  },
  delete: {
    signature: 'func delete(m map[K]V, key K)',
    description: 'Удаляет ключ из map.',
  },
  panic: {
    signature: 'func panic(v any)',
    description: 'Останавливает выполнение и запускает панику.',
  },
  recover: {
    signature: 'func recover() any',
    description: 'В deferred-вызове останавливает панику и возвращает её значение.',
    returns: 'any',
  },
};

function normalizeType(typeStr: string): string {
  return typeStr.replace(/\s+/g, ' ').trim();
}

function baseTypeName(typeStr: string): string {
  const t = normalizeType(typeStr);
  if (t.startsWith('[]')) return 'slice';
  if (t.startsWith('map[')) return 'map';
  if (t.startsWith('chan ')) return 'chan';
  if (t.startsWith('*')) return 'pointer';
  if (t.startsWith('struct{') || t.startsWith('struct {')) return 'struct';
  const word = t.split(/[\[\]\s*]/)[0];
  return word || t;
}

export function lookupType(typeStr: string, lang: Language): TypeInfo | null {
  const t = normalizeType(typeStr);
  const dict = lang === 'ru' ? typesRu : typesEn;
  const base = baseTypeName(t);

  if (dict[base]) return dict[base];

  if (t.startsWith('[]')) {
    return lang === 'ru'
      ? { title: t, description: 'Слайс — динамический массив ссылок на элементы одного типа.' }
      : { title: t, description: 'Slice — dynamic array of elements of one type.' };
  }
  if (t.startsWith('map[')) {
    return lang === 'ru'
      ? { title: t, description: 'Map — ассоциативный массив ключ → значение.' }
      : { title: t, description: 'Map — hash table from keys to values.' };
  }
  if (t.startsWith('chan ')) {
    return lang === 'ru'
      ? { title: t, description: 'Channel — тип для обмена данными между goroutine.' }
      : { title: t, description: 'Channel — typed conduit for goroutine communication.' };
  }
  if (t.startsWith('*')) {
    return lang === 'ru'
      ? { title: t, description: 'Указатель на значение типа ' + t.slice(1) + '.' }
      : { title: t, description: 'Pointer to a value of type ' + t.slice(1) + '.' };
  }
  if (t.startsWith('struct')) {
    return lang === 'ru'
      ? { title: t, description: 'Struct — составной тип из именованных полей.' }
      : { title: t, description: 'Struct — composite type with named fields.' };
  }
  if (t.startsWith('interface')) {
    return lang === 'ru'
      ? { title: t, description: 'Interface — набор методов; любой тип с этими методами подходит.' }
      : { title: t, description: 'Interface — method set; any type implementing them satisfies it.' };
  }
  return null;
}

export function formatTypeTip(info: TypeInfo | null, typeStr: string): string {
  if (!info) return typeStr;
  return `${info.title}\n${info.description}`;
}

export function lookupFunc(name: string, lang: Language): FuncInfo | null {
  const dict = lang === 'ru' ? stdlibRu : stdlibEn;
  return dict[name] ?? null;
}

export function formatFuncTip(info: FuncInfo, name: string): string {
  let tip = `${name}\n${info.signature}\n${info.description}`;
  if (info.returns) tip += `\n→ ${info.returns}`;
  return tip;
}
