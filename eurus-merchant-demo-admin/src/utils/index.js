
import { createHash } from 'crypto';

export function computeSHA256(inputStr) {
  const hash = createHash('sha256').update(inputStr).digest('hex');
  return hash; 
}

export function uuidv4() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    var r = Math.random() * 16 | 0
    var v = c == 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}
