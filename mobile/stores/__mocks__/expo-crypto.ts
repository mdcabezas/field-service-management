// Mock for expo-crypto
export const randomUUID = () => 'mock-uuid-' + Math.random().toString(36).substr(2, 9);
export const getRandomValues = <T extends Uint8Array>(array: T): T => {
  for (let i = 0; i < array.length; i++) {
    array[i] = Math.floor(Math.random() * 256);
  }
  return array;
};
export const digestStringAsync = async () => 'mock-digest';
export const digest = async () => new ArrayBuffer(32);

export default {
  randomUUID,
  getRandomValues,
  digestStringAsync,
  digest,
};