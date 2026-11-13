'use strict';
const $ = id => document.getElementById(id);

/* ---------- i18n ---------- */
const L = {
  ru: {
    back: 'Назад', forward: 'Вперёд', up: 'Вверх', refresh: 'Обновить', filter: 'Фильтр',
    newFolder: 'Создать папку', cut: 'Вырезать', copy: 'Копировать', paste: 'Вставить', rename: 'Переименовать',
    delete: 'Удалить', download: 'Скачать', open: 'Открыть', selectAll: 'Выделить всё',
    colName: 'Имя', colDate: 'Дата изменения', colType: 'Тип', colSize: 'Размер',
    computer: 'Этот компьютер', disks: 'Диски', freeOf: '{0} свободно из {1}',
    folder: 'Папка с файлами', file: 'Файл', fileExt: 'Файл "{0}"', link: 'Ссылка',
    emptyDir: 'Эта папка пуста.', noMatch: 'Нет элементов, удовлетворяющих фильтру.',
    items: ['элемент', 'элемента', 'элементов'], selected: 'Выбрано: {0}', files: 'файлов',
    clipCopy: 'В буфере (копирование): {0}', clipMove: 'В буфере (перемещение): {0}',
    ops: 'Операции', opsEmpty: 'Нет операций. Скопируйте или переместите файлы — операции появятся здесь и будут выполняться по очереди.',
    opsActive: 'выполняется: {0}, в очереди: {1}', clearDone: 'Очистить завершённые',
    opCopy: 'Копирование', opMove: 'Перемещение', opDelete: 'Удаление',
    queued: 'В очереди', scan: 'Подсчёт…', sync: 'Сброс на диск…', done: 'Готово', doneErr: 'Готово, ошибок: {0}',
    failed: 'Ошибка', canceled: 'Отменено', cancel: 'Отмена', remove: 'Убрать',
    of: '{0} из {1}', speed: 'Скорость', left: 'Осталось', elapsed: 'Прошло', total: 'Всего', took: 'Время', avg: 'Средняя скорость',
    errors: 'Ошибки ({0})', sec: 'с', min: 'мин', hour: 'ч', perSec: '/с',
    units: ['Б', 'КБ', 'МБ', 'ГБ', 'ТБ', 'ПБ'],
    newFolderTitle: 'Новая папка', newFolderName: 'Новая папка', renameTitle: 'Переименование', ok: 'ОК',
    delTitle: 'Удаление', delText: 'Удалить безвозвратно {0}?', delOne: '«{0}»',
    conflictTitle: 'Замена или пропуск файлов', conflictText: 'В папке назначения «{1}» уже есть элементы с такими именами ({0}):',
    andMore: '…и ещё {0}', replace: 'Заменить', skip: 'Пропустить', keepBoth: 'Сохранить оба',
    size: 'Размер', sizeTitle: 'Размер: {0}', sizeTotal: 'Размер', filesN: 'Файлов', foldersN: 'Папок', bytes: 'байт',
    counting: 'Идёт подсчёт…', unreadable: 'Не удалось прочитать элементов: {0}', close: 'Закрыть',
    gone: 'Папка «{0}» больше не существует.', offline: 'Нет связи с сервером',
    sameFolder: 'Конечная папка совпадает с исходной.', intoSelf: 'Нельзя переместить папку в саму себя.',
  },
  en: {
    back: 'Back', forward: 'Forward', up: 'Up', refresh: 'Refresh', filter: 'Filter',
    newFolder: 'New folder', cut: 'Cut', copy: 'Copy', paste: 'Paste', rename: 'Rename',
    delete: 'Delete', download: 'Download', open: 'Open', selectAll: 'Select all',
    colName: 'Name', colDate: 'Date modified', colType: 'Type', colSize: 'Size',
    computer: 'This PC', disks: 'Disks', freeOf: '{0} free of {1}',
    folder: 'File folder', file: 'File', fileExt: '{0} file', link: 'Link',
    emptyDir: 'This folder is empty.', noMatch: 'No items match the filter.',
    items: ['item', 'items', 'items'], selected: '{0} selected', files: 'files',
    clipCopy: 'Clipboard (copy): {0}', clipMove: 'Clipboard (move): {0}',
    ops: 'Operations', opsEmpty: 'No operations. Copy or move files — operations appear here and run one after another.',
    opsActive: 'running: {0}, queued: {1}', clearDone: 'Clear finished',
    opCopy: 'Copying', opMove: 'Moving', opDelete: 'Deleting',
    queued: 'Queued', scan: 'Calculating…', sync: 'Flushing to disk…', done: 'Done', doneErr: 'Done, errors: {0}',
    failed: 'Failed', canceled: 'Canceled', cancel: 'Cancel', remove: 'Dismiss',
    of: '{0} of {1}', speed: 'Speed', left: 'Time left', elapsed: 'Elapsed', total: 'Total', took: 'Time', avg: 'Average speed',
    errors: 'Errors ({0})', sec: 's', min: 'min', hour: 'h', perSec: '/s',
    units: ['B', 'KB', 'MB', 'GB', 'TB', 'PB'],
    newFolderTitle: 'New folder', newFolderName: 'New folder', renameTitle: 'Rename', ok: 'OK',
    delTitle: 'Delete', delText: 'Permanently delete {0}?', delOne: '"{0}"',
    conflictTitle: 'Replace or skip files', conflictText: 'The destination "{1}" already has items with these names ({0}):',
    andMore: '…and {0} more', replace: 'Replace', skip: 'Skip', keepBoth: 'Keep both',
    size: 'Size', sizeTitle: 'Size: {0}', sizeTotal: 'Size', filesN: 'Files', foldersN: 'Folders', bytes: 'bytes',
    counting: 'Calculating…', unreadable: 'Items that could not be read: {0}', close: 'Close',
    gone: 'Folder "{0}" no longer exists.', offline: 'No connection to the server',
    sameFolder: 'The destination folder is the same as the source.', intoSelf: 'Cannot move a folder into itself.',
  },
};
let lang = localStorage.getItem('lfeLang') || ((navigator.language || 'en').toLowerCase().startsWith('ru') ? 'ru' : 'en');
const t = (k, ...a) => String(L[lang][k] ?? k).replace(/\{(\d)\}/g, (_, i) => a[i]);
function plural(n, forms) {
  if (lang !== 'ru') return n === 1 ? forms[0] : forms[1];
  const a = n % 10, b = n % 100;
  return a === 1 && b !== 11 ? forms[0] : a >= 2 && a <= 4 && (b < 12 || b > 14) ? forms[1] : forms[2];
}
const nItems = n => `${n} ${plural(n, L[lang].items)}`;

/* ---------- helpers ---------- */
const num = (v, d) => v.toLocaleString(lang, { minimumFractionDigits: d, maximumFractionDigits: d });
function fmtSize(b) {
  let i = 0, v = b;
  while (v >= 1024 && i < 5) { v /= 1024; i++; }
  return num(v, i === 0 ? 0 : v < 10 ? 2 : v < 100 ? 1 : 0) + ' ' + L[lang].units[i];
}
function fmtDur(s) {
  s = Math.max(0, Math.round(s));
  if (s < 60) return `${s} ${t('sec')}`;
  if (s < 3600) return `${Math.floor(s / 60)} ${t('min')} ${s % 60} ${t('sec')}`;
  return `${Math.floor(s / 3600)} ${t('hour')} ${Math.floor(s % 3600 / 60)} ${t('min')}`;
}
const fmtDate = ms => new Date(ms).toLocaleString(lang, { dateStyle: 'short', timeStyle: 'short' });
const join = (p, n) => (p === '/' ? '/' : p + '/') + n;
const parent = p => p.slice(0, p.lastIndexOf('/')) || '/';
const base = p => p.slice(p.lastIndexOf('/') + 1);
const el = (tag, cls, text) => { const e = document.createElement(tag); if (cls) e.className = cls; if (text != null) e.textContent = text; return e; };
const ICON = {
  folder: '<svg class="fi folder" viewBox="0 0 24 24"><path d="M2.5 6.5a2 2 0 0 1 2-2h4.6l2 2h8.4a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-15a2 2 0 0 1-2-2z"/></svg>',
  file: '<svg class="fi file" viewBox="0 0 24 24"><path d="M6 3h8l5 5v12a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z"/><path d="M14 3v5h5"/></svg>',
  disk: '<svg class="fi disk" viewBox="0 0 24 24"><rect x="3" y="13" width="18" height="6" rx="1.5"/><path d="M5 13l2-7h10l2 7M16.5 16h.5"/></svg>',
  pc: '<svg class="fi disk" viewBox="0 0 24 24"><rect x="3" y="4" width="18" height="12" rx="1.5"/><path d="M8 20h8M12 16v4"/></svg>',
};
const icon = k => { const s = el('span'); s.innerHTML = ICON[k]; return s.firstChild; };

async function api(method, url, body) {
  const r = await fetch(url, {
    method, headers: { 'Content-Type': 'application/json', 'X-LFE': '1' },
    body: body ? JSON.stringify(body) : undefined,
  });
  const j = await r.json().catch(() => ({}));
  if (!r.ok) throw Object.assign(new Error(j.error || r.statusText || 'HTTP ' + r.status), { status: r.status });
  return j;
}
function toast(msg) {
  const d = el('div', 'toast', msg);
  $('toasts').append(d);
  setTimeout(() => d.remove(), 6000);
}

/* ---------- state ---------- */
const S = {
  path: '/', entries: [], view: [], sel: new Set(), anchor: 0, sort: { key: 'name', desc: false },
  dirSizes: new Map(), measuring: new Set(), sig: null, loadErr: '', lastLoad: 0, lastDisks: 0, offline: false,
  filter: '', clip: null, disks: [], jobs: [], jobEls: new Map(), prev: new Map(), order: '', drag: null, seq: 0, tick: 0,
};
const TEMP = /^\.lfe-\d+-\d+\.part$/;
const selPaths = () => S.view.filter(e => S.sel.has(e.name)).map(e => join(S.path, e.name));

/* ---------- navigation ---------- */
function hashPath() {
  try { return '/' + location.hash.slice(1).split('/').filter(Boolean).map(decodeURIComponent).join('/'); }
  catch { return '/'; }
}
function go(p) {
  const h = '#' + p.split('/').map(encodeURIComponent).join('/');
  if (hashPath() === p) load(); else location.hash = h;
}
window.addEventListener('hashchange', () => { S.sel.clear(); S.filter = $('filter').value = ''; load(); });

// load reads the current folder. keep=true is a background refresh: it never
// shows errors and does not touch the DOM when nothing has changed, so it can
// run at any moment without disturbing clicks, selection or scrolling.
async function load(keep) {
  const p = hashPath(), seq = ++S.seq;
  let path = p, entries = [], err = '';
  S.lastLoad = Date.now();
  try {
    const r = await api('GET', '/api/list?path=' + encodeURIComponent(p));
    if (seq !== S.seq) return;
    path = r.path; entries = r.entries.filter(e => !TEMP.test(e.name));
  } catch (e) {
    if (seq !== S.seq) return;
    if (e.status === 404 && p !== '/') { // folder was deleted, moved or its disk unplugged
      toast(t('gone', p)); S.sig = null; go(parent(p)); return;
    }
    if (keep) return; // temporary failure: keep what is on screen
    toast(e.message); err = e.message;
  }
  const sig = path + '\n' + JSON.stringify(entries);
  if (keep && sig === S.sig) return;
  S.sig = sig; S.path = path; S.entries = entries; S.loadErr = err;
  render();
}

/* ---------- file list ---------- */
function typeOf(e) {
  if (e.dir) return t('folder');
  const i = e.name.lastIndexOf('.');
  return i > 0 && i < e.name.length - 1 && e.name.length - i <= 9 ? t('fileExt', e.name.slice(i + 1).toUpperCase()) : t('file');
}
const sz = e => e.dir ? S.dirSizes.get(join(S.path, e.name)) ?? -1 : e.size; // measured folder size, if any
function sortEntries(list) {
  const { key, desc } = S.sort, col = new Intl.Collator(lang, { numeric: true, sensitivity: 'base' });
  const cmp = {
    name: (a, b) => col.compare(a.name, b.name),
    size: (a, b) => sz(a) - sz(b) || col.compare(a.name, b.name),
    mtime: (a, b) => a.mtime - b.mtime,
    type: (a, b) => col.compare(typeOf(a), typeOf(b)) || col.compare(a.name, b.name),
  }[key];
  return list.sort((a, b) => (b.dir - a.dir) || (desc ? -cmp(a, b) : cmp(a, b)));
}
function render() {
  const f = S.filter.toLowerCase();
  S.view = sortEntries(S.entries.filter(e => !f || e.name.toLowerCase().includes(f)));
  const names = new Set(S.view.map(e => e.name));
  for (const n of [...S.sel]) if (!names.has(n)) S.sel.delete(n);
  const cut = S.clip && S.clip.op === 'move' ? new Set(S.clip.paths) : null;

  const frag = document.createDocumentFragment();
  S.view.forEach((e, i) => {
    const tr = el('tr', 'row');
    tr.dataset.i = i; tr.draggable = true;
    if (S.sel.has(e.name)) tr.classList.add('sel');
    if (cut && cut.has(join(S.path, e.name))) tr.classList.add('cut');
    const nm = el('div', 'nm');
    nm.append(icon(e.dir ? 'folder' : 'file'), el('span', null, e.name));
    const td = el('td'); td.append(nm); td.title = e.name;
    tr.append(td, el('td', 'meta', fmtDate(e.mtime)), el('td', 'meta type', e.link ? `${typeOf(e)} (${t('link')})` : typeOf(e)),
      el('td', 'meta size', e.dir ? (S.measuring.has(join(S.path, e.name)) ? '…' : sz(e) >= 0 ? fmtSize(sz(e)) : '') : fmtSize(e.size)));
    frag.append(tr);
  });
  $('rows').replaceChildren(frag);
  $('empty').hidden = S.view.length > 0;
  $('empty').textContent = S.loadErr || (S.entries.length ? t('noMatch') : t('emptyDir'));

  document.querySelectorAll('.list th').forEach(th => {
    th.classList.toggle('sorted', th.dataset.sort === S.sort.key);
    th.classList.toggle('desc', th.dataset.sort === S.sort.key && S.sort.desc);
  });
  renderCrumbs(); renderNav(); updateUI();
  document.title = (S.path === '/' ? t('computer') : base(S.path)) + ' — Linux File Explorer';
}
function updateSel() {
  $('rows').querySelectorAll('tr.row').forEach(tr => tr.classList.toggle('sel', S.sel.has(S.view[tr.dataset.i].name)));
  updateUI();
}
function updateUI() {
  const n = S.sel.size, one = n === 1 ? S.view.find(e => S.sel.has(e.name)) : null;
  $('tbCut').disabled = $('tbCopy').disabled = $('tbDelete').disabled = n === 0;
  $('tbRename').disabled = n !== 1;
  $('tbDownload').disabled = !one || one.dir;
  $('tbPaste').disabled = !S.clip;
  $('btnUp').disabled = S.path === '/';
  $('stCount').textContent = nItems(S.view.length);
  const size = S.view.filter(e => S.sel.has(e.name) && !e.dir).reduce((a, e) => a + e.size, 0);
  $('stSel').textContent = n ? t('selected', nItems(n)) + (size ? ` · ${fmtSize(size)}` : '') : '';
  $('stClip').textContent = S.clip ? t(S.clip.op === 'move' ? 'clipMove' : 'clipCopy', nItems(S.clip.paths.length)) : '';
}
function renderCrumbs() {
  const c = $('crumbs'); c.replaceChildren();
  const parts = S.path.split('/').filter(Boolean);
  let acc = '';
  const add = (label, p) => { const b = el('button', null, label); b.dataset.path = p; b.onclick = () => go(p); c.append(b); };
  add(t('computer'), '/');
  for (const part of parts) { acc += '/' + part; c.append(el('i', null, '›')); add(part, acc); }
  c.scrollLeft = c.scrollWidth;
}
function renderNav() {
  const nav = $('nav'); nav.replaceChildren();
  const item = (ic, name, p, d) => {
    const it = el('div', 'item'); it.dataset.path = p;
    if (p === '/' ? S.path === '/' : S.path === p || S.path.startsWith(p + '/')) it.classList.add('active');
    const info = el('div', 'info'); info.append(el('div', 'name', name));
    if (d && d.total) {
      const used = 1 - d.free / d.total, bar = el('div', 'bar' + (used > 0.9 ? ' full' : '')), b = el('b');
      b.style.width = (used * 100).toFixed(1) + '%'; bar.append(b);
      info.append(bar, el('div', 'sub', t('freeOf', fmtSize(d.free), fmtSize(d.total))));
      it.title = [d.device, d.fs].filter(Boolean).join(' · ');
    }
    it.append(icon(ic), info); it.onclick = () => go(p);
    nav.append(it);
  };
  item('pc', t('computer'), '/');
  if (S.disks.length) nav.append(el('h4', null, t('disks')));
  for (const d of S.disks) item('disk', d.name, d.path, d);
}
async function loadDisks() {
  S.lastDisks = Date.now();
  try {
    const r = await api('GET', '/api/disks');
    S.disks = r.disks; $('stVer').textContent = 'v' + r.version; renderNav();
  } catch { /* shown on next refresh */ }
}

/* ---------- selection & mouse ---------- */
const rows = $('rows'), pane = $('pane');
const rowOf = ev => ev.target.closest && ev.target.closest('tr.row');
pane.addEventListener('click', ev => {
  const tr = rowOf(ev);
  if (!tr) { if (!ev.target.closest('thead')) { S.sel.clear(); updateSel(); } return; }
  const i = +tr.dataset.i, name = S.view[i].name;
  if (ev.shiftKey) {
    const [a, b] = [Math.min(S.anchor, i), Math.max(S.anchor, i)];
    if (!ev.ctrlKey && !ev.metaKey) S.sel.clear();
    for (let k = a; k <= b; k++) S.sel.add(S.view[k].name);
  } else if (ev.ctrlKey || ev.metaKey) {
    S.sel.has(name) ? S.sel.delete(name) : S.sel.add(name); S.anchor = i;
  } else { S.sel.clear(); S.sel.add(name); S.anchor = i; }
  updateSel();
});
rows.addEventListener('dblclick', ev => { const tr = rowOf(ev); if (tr) openEntry(S.view[tr.dataset.i]); });
function openEntry(e) { if (e && e.dir) go(join(S.path, e.name)); }
document.querySelector('.list thead').addEventListener('click', ev => {
  const th = ev.target.closest('th'); if (!th) return;
  S.sort = { key: th.dataset.sort, desc: S.sort.key === th.dataset.sort ? !S.sort.desc : false };
  render();
});

/* ---------- drag & drop ---------- */
rows.addEventListener('dragstart', ev => {
  const tr = rowOf(ev); if (!tr) return;
  const name = S.view[tr.dataset.i].name;
  if (!S.sel.has(name)) { S.sel.clear(); S.sel.add(name); updateSel(); }
  S.drag = selPaths();
  ev.dataTransfer.setData('text/plain', S.drag.join('\n'));
  ev.dataTransfer.effectAllowed = 'copyMove';
});
document.addEventListener('dragend', () => { S.drag = null; clearDrop(); });
let dropEl = null;
function clearDrop() { if (dropEl) dropEl.classList.remove('drop'); dropEl = null; }
function dropTarget(ev) {
  if (!S.drag) return null;
  let node = rowOf(ev), p = null;
  if (node) { const e = S.view[node.dataset.i]; if (e.dir) p = join(S.path, e.name); }
  else if ((node = ev.target.closest('.nav .item, .crumbs button'))) p = node.dataset.path;
  if (!p || S.drag.some(s => s === p || p.startsWith(s + '/') || parent(s) === p)) return null;
  return { node, path: p };
}
const dropOp = (ev, dest) => ev.ctrlKey ? 'copy' : ev.shiftKey ? 'move'
  : S.drag[0].split('/')[1] === dest.split('/')[1] && dest !== '/' ? 'move' : 'copy';
document.addEventListener('dragover', ev => {
  const tg = dropTarget(ev);
  if (!tg) { clearDrop(); return; }
  ev.preventDefault();
  ev.dataTransfer.dropEffect = dropOp(ev, tg.path);
  if (dropEl !== tg.node) { clearDrop(); dropEl = tg.node; dropEl.classList.add('drop'); }
});
document.addEventListener('drop', ev => {
  const tg = dropTarget(ev); clearDrop();
  if (!tg) return;
  ev.preventDefault();
  transfer(dropOp(ev, tg.path), S.drag, tg.path);
  S.drag = null;
});

/* ---------- actions ---------- */
async function transfer(op, sources, dest) {
  try {
    if (op === 'move' && sources.every(s => parent(s) === dest)) { toast(t('sameFolder')); return false; }
    if (sources.some(s => dest === s || dest.startsWith(s + '/'))) { toast(t('intoSelf')); return false; }
    let conflict = 'rename';
    const c = await api('POST', '/api/check', { sources, dest });
    if (c.conflicts.length) {
      const shown = c.conflicts.slice(0, 8).join('\n') + (c.conflicts.length > 8 ? '\n' + t('andMore', c.conflicts.length - 8) : '');
      const r = await dialog({
        title: t('conflictTitle'), text: t('conflictText', c.conflicts.length, dest) + '\n\n' + shown,
        buttons: [{ id: 'overwrite', label: t('replace'), primary: true }, { id: 'skip', label: t('skip') },
          { id: 'rename', label: t('keepBoth') }, { id: '', label: t('cancel') }],
      });
      if (!r || !r.btn) return false;
      conflict = r.btn;
    }
    await api('POST', '/api/jobs', { op, sources, dest, conflict });
    $('ops').classList.remove('collapsed');
    poll();
    return true;
  } catch (e) { toast(e.message); return false; }
}
function setClip(op) {
  const paths = selPaths(); if (!paths.length) return;
  S.clip = { op, paths }; render();
}
async function paste() {
  if (!S.clip) return;
  const { op, paths } = S.clip;
  if (await transfer(op, paths, S.path) && op === 'move') { S.clip = null; render(); }
}
async function newFolder() {
  let name = t('newFolderName');
  const taken = new Set(S.entries.map(e => e.name));
  for (let i = 2; taken.has(name); i++) name = `${t('newFolderName')} (${i})`;
  const r = await dialog({ title: t('newFolderTitle'), input: name, buttons: [{ id: 'ok', label: t('ok'), primary: true }, { id: '', label: t('cancel') }] });
  if (!r || !r.btn || !r.value.trim()) return;
  try {
    await api('POST', '/api/mkdir', { path: S.path, name: r.value.trim() });
    S.sel = new Set([r.value.trim()]); await load();
  } catch (e) { toast(e.message); }
}
async function rename() {
  if (S.sel.size !== 1) return;
  const e = S.view.find(x => S.sel.has(x.name));
  const dot = e.dir ? -1 : e.name.lastIndexOf('.');
  const r = await dialog({
    title: t('renameTitle'), input: e.name, selectTo: dot > 0 ? dot : e.name.length,
    buttons: [{ id: 'ok', label: t('ok'), primary: true }, { id: '', label: t('cancel') }],
  });
  const name = r && r.btn ? r.value.trim() : '';
  if (!name || name === e.name) return;
  try {
    await api('POST', '/api/rename', { path: join(S.path, e.name), name });
    S.sel = new Set([name]); await load();
  } catch (err) { toast(err.message); }
}
async function del() {
  const paths = selPaths(); if (!paths.length) return;
  const what = paths.length === 1 ? t('delOne', base(paths[0])) : nItems(paths.length);
  const r = await dialog({ title: t('delTitle'), text: t('delText', what), buttons: [{ id: 'ok', label: t('delete'), danger: true }, { id: '', label: t('cancel') }] });
  if (!r || !r.btn) return;
  try {
    await api('POST', '/api/jobs', { op: 'delete', sources: paths });
    $('ops').classList.remove('collapsed'); poll();
  } catch (e) { toast(e.message); }
}
function download() {
  const e = S.sel.size === 1 && S.view.find(x => S.sel.has(x.name));
  if (!e || e.dir) return;
  const a = el('a'); a.href = '/api/download?path=' + encodeURIComponent(join(S.path, e.name)); a.download = e.name;
  document.body.append(a); a.click(); a.remove();
}
// measure counts the size of the selection (or of the open folder when
// nothing is selected), showing running totals while the server walks the tree.
async function measure() {
  const paths = selPaths().length ? selPaths() : [S.path];
  const name = paths.length > 1 ? nItems(paths.length) : paths[0] === '/' ? t('computer') : base(paths[0]);
  const ctl = new AbortController();
  dialog({ title: t('sizeTitle', name), text: '', buttons: [{ id: 'ok', label: t('close'), primary: true }] }).then(() => ctl.abort());
  const mine = dlgClose;
  const busy = on => {
    paths.forEach(p => on ? S.measuring.add(p) : S.measuring.delete(p));
    if (dlgClose === mine) { $('dlgBusy').hidden = !on; $('dlgNote').hidden = !on; $('dlgNote').textContent = on ? t('counting') : ''; }
    render();
  };
  busy(true);
  const show = r => {
    if (dlgClose !== mine) return;
    if (!r.done) $('dlgNote').textContent = t('counting') + (r.cur ? ' ' + r.cur : '');
    $('dlgBody').textContent = `${t('sizeTotal')}: ${fmtSize(r.bytes)} (${num(r.bytes, 0)} ${t('bytes')})\n${t('filesN')}: ${num(r.files, 0)}\n${t('foldersN')}: ${num(r.dirs, 0)}`
      + (r.errors ? '\n' + t('unreadable', r.errors) : '');
  };
  try {
    const resp = await fetch('/api/size', {
      method: 'POST', headers: { 'Content-Type': 'application/json', 'X-LFE': '1' },
      body: JSON.stringify({ paths }), signal: ctl.signal,
    });
    if (!resp.ok) throw new Error((await resp.json().catch(() => ({}))).error || 'HTTP ' + resp.status);
    const rd = resp.body.getReader(), dec = new TextDecoder();
    let buf = '';
    for (;;) {
      const { value, done } = await rd.read();
      if (done) break;
      buf += dec.decode(value, { stream: true });
      const lines = buf.split('\n'); buf = lines.pop();
      const last = lines.filter(Boolean).pop();
      if (!last) continue;
      const r = JSON.parse(last);
      show(r);
      if (r.done) paths.forEach((p, i) => S.dirSizes.set(p, r.sizes[i]));
    }
  } catch (e) {
    if (e.name !== 'AbortError') { toast(e.message); if (dlgClose === mine) mine(''); }
  }
  busy(false);
}
function selectAll() { S.view.forEach(e => S.sel.add(e.name)); updateSel(); }

/* ---------- dialog ---------- */
let dlgClose = null;
function dialog({ title, text, input, selectTo, buttons }) {
  return new Promise(resolve => {
    const inp = $('dlgInput');
    $('dlgTitle').textContent = title;
    $('dlgBody').textContent = text || '';
    $('dlgBusy').hidden = $('dlgNote').hidden = true;
    inp.hidden = input == null; inp.value = input ?? '';
    const close = btn => { $('modal').hidden = true; dlgClose = null; resolve({ btn, value: inp.value }); pane.focus(); };
    dlgClose = close;
    const bs = $('dlgBtns'); bs.replaceChildren();
    let first = null;
    for (const b of buttons) {
      const e = el('button', b.primary ? 'primary' : b.danger ? 'danger' : '', b.label);
      e.onclick = () => close(b.id); bs.append(e);
      if (!first && b.id) first = e;
    }
    $('modal').hidden = false;
    if (input != null) { inp.focus(); inp.setSelectionRange(0, selectTo ?? inp.value.length); } else first.focus();
    inp.onkeydown = ev => { if (ev.key === 'Enter') close(buttons[0].id); };
  });
}
$('modal').addEventListener('mousedown', ev => { if (ev.target === $('modal') && dlgClose) dlgClose(''); });

/* ---------- context menu ---------- */
const menu = $('menu');
function showMenu(x, y, items) {
  menu.replaceChildren();
  for (const it of items) {
    if (!it) { menu.append(el('hr')); continue; }
    const b = el('button'); b.append(el('span', null, it.label)); if (it.key) b.append(el('kbd', null, it.key));
    b.disabled = !!it.disabled; b.onclick = () => { hideMenu(); it.run(); };
    menu.append(b);
  }
  menu.hidden = false;
  menu.style.left = Math.min(x, innerWidth - menu.offsetWidth - 4) + 'px';
  menu.style.top = Math.min(y, innerHeight - menu.offsetHeight - 4) + 'px';
}
const hideMenu = () => { menu.hidden = true; };
document.addEventListener('mousedown', ev => { if (!menu.contains(ev.target)) hideMenu(); });
window.addEventListener('blur', hideMenu);
pane.addEventListener('contextmenu', ev => {
  ev.preventDefault();
  const tr = rowOf(ev);
  if (tr) {
    const e = S.view[tr.dataset.i];
    if (!S.sel.has(e.name)) { S.sel.clear(); S.sel.add(e.name); S.anchor = +tr.dataset.i; updateSel(); }
    const one = S.sel.size === 1;
    showMenu(ev.clientX, ev.clientY, [
      one && e.dir ? { label: t('open'), key: 'Enter', run: () => openEntry(e) } : { label: t('download'), disabled: !one, run: download },
      null,
      { label: t('cut'), key: 'Ctrl+X', run: () => setClip('move') },
      { label: t('copy'), key: 'Ctrl+C', run: () => setClip('copy') },
      { label: t('paste'), key: 'Ctrl+V', disabled: !S.clip, run: paste },
      null,
      { label: t('rename'), key: 'F2', disabled: !one, run: rename },
      { label: t('delete'), key: 'Del', run: del },
      null,
      { label: t('size'), run: measure },
    ]);
  } else {
    S.sel.clear(); updateSel();
    showMenu(ev.clientX, ev.clientY, [
      { label: t('paste'), key: 'Ctrl+V', disabled: !S.clip, run: paste },
      { label: t('newFolder'), run: newFolder },
      null,
      { label: t('selectAll'), key: 'Ctrl+A', run: selectAll },
      { label: t('refresh'), run: () => { load(); loadDisks(); } },
      { label: t('size'), run: measure },
    ]);
  }
});

/* ---------- keyboard ---------- */
document.addEventListener('keydown', ev => {
  if (dlgClose) { if (ev.key === 'Escape') dlgClose(''); return; }
  if (ev.key === 'Escape') { hideMenu(); if (ev.target === $('filter')) { $('filter').value = S.filter = ''; render(); pane.focus(); } else { S.sel.clear(); updateSel(); } return; }
  if (ev.target.tagName === 'INPUT') return;
  const ctrl = ev.ctrlKey || ev.metaKey, k = ev.key.toLowerCase(), code = ev.code;
  const done = () => ev.preventDefault();
  const is = c => code === 'Key' + c || (!code && k === c.toLowerCase()); // code: layout-independent (works with a Russian layout)
  if (ctrl && is('A')) { selectAll(); done(); }
  else if (ctrl && is('C')) { setClip('copy'); done(); }
  else if (ctrl && is('X')) { setClip('move'); done(); }
  else if (ctrl && is('V')) { paste(); done(); }
  else if (ctrl && is('F')) { $('filter').focus(); done(); }
  else if (k === 'delete') { del(); done(); }
  else if (k === 'f2') { rename(); done(); }
  else if (k === 'backspace') { if (S.path !== '/') go(parent(S.path)); done(); }
  else if (k === 'enter') { if (S.sel.size === 1) openEntry(S.view.find(e => S.sel.has(e.name))); done(); }
  else if (['arrowdown', 'arrowup', 'home', 'end'].includes(k) && S.view.length) {
    let i = S.sel.size ? S.anchor : -1;
    i = k === 'home' ? 0 : k === 'end' ? S.view.length - 1 : Math.max(0, Math.min(S.view.length - 1, i + (k === 'arrowdown' ? 1 : -1)));
    S.sel.clear(); S.sel.add(S.view[i].name); S.anchor = i; updateSel();
    rows.children[i].scrollIntoView({ block: 'nearest' }); done();
  }
});

/* ---------- operations panel ---------- */
const OPN = { copy: 'opCopy', move: 'opMove', delete: 'opDelete' };
const active = j => j.state === 'queued' || j.state === 'running';
function jobCard(j) {
  const c = el('div', 'job'), r = {};
  const top = el('div', 'top');
  r.title = el('div', 'title'); r.badge = el('span', 'badge'); r.act = el('button', 'act');
  top.append(r.title, r.badge, r.act);
  const prog = el('div', 'prog'), track = el('div', 'track');
  r.bar = el('b'); track.append(r.bar); r.pct = el('span', 'pct'); prog.append(track, r.pct);
  r.stats = el('div', 'stats'); r.cur = el('div', 'cur');
  r.errs = el('details', 'errs'); r.errSum = el('summary'); r.errBody = el('div'); r.errs.append(r.errSum, r.errBody);
  c.append(top, prog, r.stats, r.cur, r.errs);
  const what = j.sources.length === 1 ? base(j.sources[0]) : nItems(j.sources.length);
  r.title.append(el('b', null, t(OPN[j.op]) + ': '), what + (j.op === 'delete' ? ` (${parent(j.sources[0])})` : ` → ${j.dest}`));
  r.title.title = j.sources.join('\n') + (j.dest ? '\n→ ' + j.dest : '');
  r.act.onclick = async () => {
    try { await api(active(c._j) ? 'POST' : 'DELETE', `/api/jobs/${j.id}` + (active(c._j) ? '/cancel' : '')); } catch (e) { toast(e.message); }
    poll();
  };
  c._r = r;
  return c;
}
function stat(parent_, label, value) {
  const s = el('span', null, label ? label + ': ' : ''); s.append(el('b', null, value)); parent_.append(s);
}
function updateJob(c, j) {
  const r = c._r; c._j = j;
  const run = j.state === 'running', fin = !active(j);
  const byBytes = j.totalBytes > 0;
  let frac = byBytes ? j.doneBytes / j.totalBytes : j.totalFiles > 0 ? j.doneFiles / j.totalFiles : 0;
  frac = j.state === 'done' ? 1 : Math.max(0, Math.min(1, frac));
  c.className = 'job ' + j.state + (j.errCount ? ' haserr' : '') + (run && (j.phase === 'scan' || j.phase === 'sync') ? ' indet' : '');
  r.badge.textContent = run ? (j.phase === 'scan' || j.phase === 'sync' ? t(j.phase) : t(OPN[j.op]))
    : j.state === 'done' ? (j.errCount ? t('doneErr', j.errCount) : t('done')) : t(j.state);
  r.act.textContent = fin ? t('remove') : t('cancel');
  r.bar.style.width = (frac * 100).toFixed(1) + '%';
  r.pct.textContent = j.state === 'queued' ? '' : Math.floor(frac * 100) + '%';

  r.stats.replaceChildren();
  const filesStr = `${num(j.doneFiles, 0)} / ${num(j.totalFiles, 0)}`;
  if (run) {
    if (byBytes) stat(r.stats, '', t('of', fmtSize(j.doneBytes), fmtSize(j.totalBytes)));
    if (j.totalFiles) stat(r.stats, t('files'), filesStr);
    if (j.op !== 'delete' && j.phase !== 'scan') {
      stat(r.stats, t('speed'), fmtSize(j.speed) + t('perSec'));
      if (j.speed > 0 && byBytes) stat(r.stats, t('left'), fmtDur((j.totalBytes - j.doneBytes) / j.speed));
    }
    stat(r.stats, t('elapsed'), fmtDur(j.elapsed / 1000));
  } else if (fin && j.started) {
    const secs = Math.max(j.elapsed / 1000, 0.001);
    if (byBytes) stat(r.stats, t('total'), j.state === 'done' ? fmtSize(j.totalBytes) : t('of', fmtSize(j.doneBytes), fmtSize(j.totalBytes)));
    if (j.totalFiles) stat(r.stats, t('files'), filesStr);
    stat(r.stats, t('took'), fmtDur(secs));
    if (byBytes && j.op !== 'delete' && secs >= 1) stat(r.stats, t('avg'), fmtSize(j.doneBytes / secs) + t('perSec'));
  }
  r.cur.textContent = run ? j.current : '';
  r.cur.hidden = !run || !j.current;
  r.errs.hidden = !j.errCount;
  if (j.errCount) { r.errSum.textContent = t('errors', j.errCount); r.errBody.textContent = j.errors.join('\n'); }
}
function renderJobs() {
  const rank = j => j.state === 'running' ? 0 : j.state === 'queued' ? 1 : 2;
  const list = [...S.jobs].sort((a, b) => rank(a) - rank(b) || (rank(a) === 2 ? b.finished - a.finished : a.id - b.id));
  const ids = new Set(list.map(j => j.id));
  for (const [id, c] of S.jobEls) if (!ids.has(id)) { c.remove(); S.jobEls.delete(id); }
  for (const j of list) {
    let c = S.jobEls.get(j.id);
    if (!c) { c = jobCard(j); S.jobEls.set(j.id, c); }
    updateJob(c, j);
  }
  const order = list.map(j => j.id).join();
  if (order !== S.order) { S.order = order; $('jobs').append(...list.map(j => S.jobEls.get(j.id))); }
  const nRun = list.filter(j => j.state === 'running').length, nQ = list.filter(j => j.state === 'queued').length;
  $('opsSummary').textContent = nRun + nQ ? t('opsActive', nRun, nQ) : '';
  $('opsClear').hidden = !list.some(j => !active(j));
}
function setOffline(off) {
  S.offline = off;
  $('stNet').textContent = off ? t('offline') : '';
}
// Returning to the tab: show the current state at once instead of a stale one.
const wake = () => { if (!document.hidden) { poll(); if (!dlgClose && !S.drag) load(true); loadDisks(); } };
document.addEventListener('visibilitychange', wake);
window.addEventListener('focus', wake);
let pollTimer = 0, pollGen = 0, polled = false;
async function poll() {
  clearTimeout(pollTimer);
  const gen = ++pollGen;
  try {
    const jobs = (await api('GET', '/api/jobs')).jobs;
    if (gen !== pollGen) return; // a newer poll is already running
    S.jobs = jobs;
    let finishedNow = false;
    for (const j of S.jobs) {
      // A job counts as "just finished" when we saw it active before, or when
      // it is new to us and already done (instant moves finish before the
      // first poll ever sees them running).
      const before = S.prev.get(j.id);
      if (!active(j) && (before ? active({ state: before }) : polled)) finishedNow = true;
      S.prev.set(j.id, j.state);
    }
    polled = true;
    if (finishedNow && S.dirSizes.size) { S.dirSizes.clear(); S.sig = null; } // sizes are stale after any operation
    renderJobs();
    const run = S.jobs.find(j => j.state === 'running');
    S.tick++;
    const touches = run && (run.dest === S.path || run.sources.some(s => parent(s) === S.path));
    // Background refresh: often while a job writes into the open folder,
    // rarely otherwise (to pick up changes made by other tabs or programs).
    const now = Date.now(), wasOffline = S.offline, quiet = !S.drag && !dlgClose && menu.hidden && !document.hidden;
    setOffline(false);
    if (finishedNow || wasOffline || (quiet && now - S.lastLoad > (touches ? 4000 : 15000))) load(true);
    if (finishedNow || wasOffline || (!document.hidden && now - S.lastDisks > (run ? 10000 : 30000))) loadDisks();
  } catch { if (gen === pollGen) setOffline(true); }
  if (gen !== pollGen) return;
  pollTimer = setTimeout(poll, S.jobs.some(active) ? 1000 : 4000);
}

/* ---------- wiring ---------- */
function applyLang() {
  document.documentElement.lang = lang;
  document.querySelectorAll('[data-t]').forEach(e => { e.textContent = t(e.dataset.t); });
  document.querySelectorAll('[data-tip]').forEach(e => { e.title = t(e.dataset.tip); });
  const th = document.querySelectorAll('.list th');
  ['colName', 'colDate', 'colType', 'colSize'].forEach((k, i) => { th[i].textContent = t(k); });
  $('filter').placeholder = t('filter');
  $('opsTitle').textContent = t('ops');
  $('opsClear').textContent = t('clearDone');
  $('jobs').dataset.empty = t('opsEmpty');
  $('btnLang').textContent = lang === 'ru' ? 'English' : 'Русский';
  for (const c of S.jobEls.values()) c.remove();
  S.jobEls.clear(); S.order = '';
  render(); renderJobs();
}
$('btnLang').onclick = () => { lang = lang === 'ru' ? 'en' : 'ru'; localStorage.setItem('lfeLang', lang); applyLang(); };
$('btnBack').onclick = () => history.back();
$('btnFwd').onclick = () => history.forward();
$('btnUp').onclick = () => go(parent(S.path));
$('btnRefresh').onclick = () => { load(); loadDisks(); };
$('tbNew').onclick = newFolder;
$('tbCut').onclick = () => setClip('move');
$('tbCopy').onclick = () => setClip('copy');
$('tbPaste').onclick = paste;
$('tbRename').onclick = rename;
$('tbDelete').onclick = del;
$('tbDownload').onclick = download;
$('tbSize').onclick = measure;
$('filter').oninput = ev => { S.filter = ev.target.value; render(); };
$('opsToggle').onclick = () => $('ops').classList.toggle('collapsed');
$('opsClear').onclick = async () => { try { await api('POST', '/api/jobs/clear'); } catch (e) { toast(e.message); } poll(); };

applyLang();
load(); loadDisks(); poll();
