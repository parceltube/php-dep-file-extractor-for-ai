// ============================================================
// PHP Dependency Extractor - Frontend  v0.1.0
// ============================================================

const $ = (sel) => document.querySelector(sel);
const $$ = (sel) => document.querySelectorAll(sel);

// State
const state = {
    projectPath: '',
    outputPath: '',
    outputIsAuto: true,
    treeData: null,
    selectedFiles: new Set(),
    dependencies: [],
    includes: [],
    checkedIncludes: new Set(),
    searchFilter: '',
};

// ============================================================
// API helpers
// ============================================================

async function api(path, body) {
    const resp = await fetch(path, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
    });
    const data = await resp.json();
    if (data.error) throw new Error(data.error);
    return data;
}

// ============================================================
// Theme (now inside Settings modal)
// ============================================================

function applyTheme(theme) {
    if (theme === 'system') {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        document.body.setAttribute('data-theme', prefersDark ? 'dark' : 'light');
    } else {
        document.body.setAttribute('data-theme', theme);
    }
    localStorage.setItem('pde-theme', theme);

    // Sync radio buttons
    const radio = document.querySelector(`input[name="theme"][value="${theme}"]`);
    if (radio) radio.checked = true;
}

document.querySelectorAll('input[name="theme"]').forEach(radio => {
    radio.addEventListener('change', (e) => applyTheme(e.target.value));
});

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    const saved = localStorage.getItem('pde-theme') || 'dark';
    if (saved === 'system') applyTheme('system');
});

(function initTheme() {
    const saved = localStorage.getItem('pde-theme') || 'dark';
    applyTheme(saved);
})();

// ============================================================
// Font size (now inside Settings modal)
// ============================================================

function applyFontSize(px) {
    document.body.style.setProperty('--base-font-size', px + 'px');
    document.body.style.fontSize = px + 'px';
    $('#fontSizeValue').textContent = px + 'px';
    $('#fontSizeRange').value = px;
    localStorage.setItem('pde-font-size', px);
}

$('#fontSizeRange').addEventListener('input', (e) => applyFontSize(parseInt(e.target.value)));

$('#fontDec').addEventListener('click', () => {
    const cur = parseInt($('#fontSizeRange').value);
    if (cur > 11) applyFontSize(cur - 1);
});

$('#fontInc').addEventListener('click', () => {
    const cur = parseInt($('#fontSizeRange').value);
    if (cur < 18) applyFontSize(cur + 1);
});

(function initFontSize() {
    const saved = parseInt(localStorage.getItem('pde-font-size')) || 13;
    applyFontSize(saved);
})();

// ============================================================
// Settings modal with tabs
// ============================================================

$('#btnSettings').addEventListener('click', () => {
    loadSettings();
    $('#settingsModal').classList.add('active');
});

$('#btnSettingsClose').addEventListener('click', () => {
    $('#settingsModal').classList.remove('active');
});

// Close on overlay click
$('#settingsModal').addEventListener('click', (e) => {
    if (e.target === e.currentTarget) {
        $('#settingsModal').classList.remove('active');
    }
});

// Tab switching
$$('.tab').forEach(tabBtn => {
    tabBtn.addEventListener('click', () => {
        $$('.tab').forEach(t => t.classList.remove('active'));
        $$('.tab-content').forEach(c => c.classList.remove('active'));
        tabBtn.classList.add('active');
        const target = tabBtn.getAttribute('data-tab');
        $(`#${target}`).classList.add('active');
    });
});

// Save mappings
$('#btnMappingsSave').addEventListener('click', async () => {
    const rows = $$('.mapping-row');
    const mappings = [];
    rows.forEach(row => {
        const inputs = row.querySelectorAll('input');
        if (inputs[0].value && inputs[1].value) {
            mappings.push({ prefix: inputs[0].value, dir: inputs[1].value });
        }
    });

    try {
        await api('/api/settings', { mappings });
        setStatus('Prefix mappings saved');
    } catch (e) {
        setStatus('Error saving settings: ' + e.message);
    }
});

$('#btnAddMapping').addEventListener('click', () => {
    addMappingRow('', '');
});

function loadSettings() {
    fetch('/api/settings')
        .then(r => r.json())
        .then(data => {
            $('#mappingsList').innerHTML = '';
            (data.mappings || []).forEach(m => addMappingRow(m.prefix, m.dir));
        });
}

function addMappingRow(prefix, dir) {
    const row = document.createElement('div');
    row.className = 'mapping-row';
    row.innerHTML = `
        <input type="text" value="${escHtml(prefix)}" placeholder="Prefix (e.g. Model_)">
        <span style="color:var(--text-dim)">&#8594;</span>
        <input type="text" value="${escHtml(dir)}" placeholder="Directory (e.g. models/)">
        <button class="btn-icon" title="Remove">&#10005;</button>
    `;
    row.querySelector('.btn-icon').addEventListener('click', () => row.remove());
    $('#mappingsList').appendChild(row);
}

// ============================================================
// Progress indicator
// ============================================================

function showProgress(title, detail) {
    $('#progressTitle').textContent = title;
    $('#progressDetail').textContent = detail || '';
    const bar = $('#progressBar');
    bar.style.width = '';
    bar.classList.add('indeterminate');
    bar.classList.remove('determinate');
    $('#progressOverlay').classList.add('active');
}

function updateProgress(detail, percent) {
    if (detail !== undefined) $('#progressDetail').textContent = detail;
    if (percent !== undefined && percent >= 0) {
        const bar = $('#progressBar');
        bar.classList.remove('indeterminate');
        bar.style.width = Math.min(100, percent) + '%';
    }
}

function hideProgress() {
    $('#progressOverlay').classList.remove('active');
}

// ============================================================
// Auto output directory
// ============================================================

function getAutoOutputPath() {
    if (!state.projectPath) return '';
    const now = new Date();
    const ts = now.getFullYear().toString() +
        String(now.getMonth() + 1).padStart(2, '0') +
        String(now.getDate()).padStart(2, '0') + '_' +
        String(now.getHours()).padStart(2, '0') +
        String(now.getMinutes()).padStart(2, '0') +
        String(now.getSeconds()).padStart(2, '0');
    return state.projectPath + '_output_' + ts;
}

function getEffectiveOutputPath() {
    if (state.outputPath && !state.outputIsAuto) return state.outputPath;
    return getAutoOutputPath();
}

function refreshOutputDisplay() {
    if (state.outputIsAuto) {
        $('#outputPath').value = '';
        $('#outputPath').placeholder = state.projectPath
            ? getAutoOutputPath()
            : 'Auto: {project}_output_{ts}';
    }
}

// ============================================================
// Toolbar actions
// ============================================================

$('#btnBrowseProject').addEventListener('click', async () => {
    setStatus('Opening folder dialog...');
    try {
        const data = await api('/api/browse', {});
        if (data.path) {
            state.projectPath = data.path;
            $('#projectPath').value = data.path;
            $('#btnScan').disabled = false;
            refreshOutputDisplay();
            setStatus('Project directory selected');
        } else {
            setStatus('No directory selected');
        }
    } catch (e) {
        setStatus('Error: ' + e.message);
    }
});

$('#btnBrowseOutput').addEventListener('click', async () => {
    setStatus('Opening folder dialog...');
    try {
        const data = await api('/api/browse', {});
        if (data.path) {
            state.outputPath = data.path;
            state.outputIsAuto = false;
            $('#outputPath').value = data.path;
            $('#outputPath').placeholder = '';
            updateCopyButton();
            setStatus('Output directory selected');
        }
    } catch (e) {
        setStatus('Error: ' + e.message);
    }
});

$('#outputPath').addEventListener('dblclick', () => {
    state.outputPath = '';
    state.outputIsAuto = true;
    $('#outputPath').value = '';
    refreshOutputDisplay();
    updateCopyButton();
    setStatus('Output reset to auto-generate');
});

$('#btnScan').addEventListener('click', async () => {
    const path = state.projectPath;
    if (!path) return;

    showProgress('Scanning Project', 'Traversing directory and building class index...');
    setStatus('Scanning project...');
    $('#btnScan').disabled = true;

    try {
        const data = await api('/api/scan', {
            path: path,
            framework: $('#framework').value,
        });

        state.treeData = data.tree;
        state.selectedFiles.clear();
        state.dependencies = [];
        state.includes = [];

        updateProgress('Building file tree...', 80);

        renderTree(data.tree);
        renderResults();
        $('#fileCount').textContent = `${data.fileCount} files`;

        updateProgress('Done!', 100);
        setTimeout(hideProgress, 400);

        setStatus(`Scanned ${data.fileCount} files, indexed ${data.indexed} classes`);
        $('#btnAnalyze').disabled = false;
    } catch (e) {
        hideProgress();
        setStatus('Scan error: ' + e.message);
    } finally {
        $('#btnScan').disabled = false;
    }
});

$('#btnAnalyze').addEventListener('click', async () => {
    if (state.selectedFiles.size === 0) {
        setStatus('No files selected');
        return;
    }

    const fileCount = state.selectedFiles.size;
    showProgress('Analyzing Dependencies', `Parsing ${fileCount} selected file${fileCount > 1 ? 's' : ''}...`);
    setStatus('Analyzing dependencies...');
    $('#btnAnalyze').disabled = true;

    try {
        const data = await api('/api/analyze', {
            files: Array.from(state.selectedFiles),
            parseIncludes: $('#parseIncludes').checked,
        });

        state.dependencies = data.dependencies || [];
        state.includes = data.includes || [];
        state.checkedIncludes.clear();

        updateProgress('Rendering results...', 90);

        renderResults();
        updateCopyButton();

        const depCount = state.dependencies.length;
        const incCount = state.includes.length;

        updateProgress(`Found ${depCount} dependencies`, 100);
        setTimeout(hideProgress, 400);

        setStatus(`Found ${depCount} dependencies` + (incCount > 0 ? `, ${incCount} includes` : ''));
    } catch (e) {
        hideProgress();
        setStatus('Analysis error: ' + e.message);
    } finally {
        $('#btnAnalyze').disabled = false;
    }
});

$('#btnCopy').addEventListener('click', async () => {
    const outputDir = getEffectiveOutputPath();
    if (!outputDir) {
        setStatus('Select a project directory first');
        return;
    }

    const files = getAllCopyFiles();
    if (files.length === 0) {
        setStatus('No files to copy');
        return;
    }

    showProgress('Copying Files', `Copying ${files.length} file${files.length > 1 ? 's' : ''} to output...`);
    setStatus('Copying files...');
    $('#btnCopy').disabled = true;

    try {
        const data = await api('/api/copy', {
            files: files,
            outputDir: outputDir,
        });

        const copied = (data.copied || []).length;
        const errors = (data.errors || []).length;

        updateProgress(`Copied ${copied} files`, 100);
        setTimeout(hideProgress, 500);

        setStatus(`Copied ${copied} files to ${outputDir}` + (errors > 0 ? ` (${errors} errors)` : ''));
    } catch (e) {
        hideProgress();
        setStatus('Copy error: ' + e.message);
    } finally {
        $('#btnCopy').disabled = false;
    }
});

// Search filter
$('#searchFilter').addEventListener('input', (e) => {
    state.searchFilter = e.target.value.toLowerCase();
    if (state.treeData) {
        renderTree(state.treeData);
    }
});

// ============================================================
// File tree rendering
// ============================================================

function renderTree(node) {
    const container = $('#treeContainer');
    container.innerHTML = '';

    if (!node || !node.children || node.children.length === 0) {
        container.innerHTML = '<div class="empty-state"><p>No PHP files found</p></div>';
        return;
    }

    const fragment = document.createDocumentFragment();
    node.children.forEach(child => {
        const el = createTreeNode(child, 0);
        if (el) fragment.appendChild(el);
    });
    container.appendChild(fragment);
}

function createTreeNode(node, depth) {
    if (state.searchFilter) {
        if (!node.isDir && !node.path.toLowerCase().includes(state.searchFilter)) {
            return null;
        }
        if (node.isDir && !hasVisibleChildren(node)) {
            return null;
        }
    }

    const div = document.createElement('div');
    div.className = 'tree-node';

    const item = document.createElement('div');
    item.className = 'tree-item';
    item.style.paddingLeft = (8 + depth * 16) + 'px';

    if (node.isDir) {
        const toggle = document.createElement('span');
        toggle.className = 'tree-toggle';
        toggle.textContent = '\u25B6';

        const icon = document.createElement('span');
        icon.className = 'tree-icon';
        icon.textContent = '\uD83D\uDCC1';

        const name = document.createElement('span');
        name.className = 'tree-name';
        name.textContent = node.name;

        item.appendChild(toggle);
        item.appendChild(icon);
        item.appendChild(name);

        const childrenDiv = document.createElement('div');
        childrenDiv.className = 'tree-children collapsed';

        let childrenLoaded = false;

        item.addEventListener('click', (e) => {
            e.stopPropagation();
            const isCollapsed = childrenDiv.classList.contains('collapsed');

            if (isCollapsed) {
                if (!childrenLoaded && node.children) {
                    node.children.forEach(child => {
                        const childEl = createTreeNode(child, depth + 1);
                        if (childEl) childrenDiv.appendChild(childEl);
                    });
                    childrenLoaded = true;
                }
                childrenDiv.classList.remove('collapsed');
                toggle.textContent = '\u25BC';
            } else {
                childrenDiv.classList.add('collapsed');
                toggle.textContent = '\u25B6';
            }
        });

        if (state.searchFilter) {
            if (node.children) {
                node.children.forEach(child => {
                    const childEl = createTreeNode(child, depth + 1);
                    if (childEl) childrenDiv.appendChild(childEl);
                });
                childrenLoaded = true;
            }
            childrenDiv.classList.remove('collapsed');
            toggle.textContent = '\u25BC';
        }

        div.appendChild(item);
        div.appendChild(childrenDiv);
    } else {
        const toggle = document.createElement('span');
        toggle.className = 'tree-toggle';
        toggle.textContent = '';

        const cb = document.createElement('input');
        cb.type = 'checkbox';
        cb.className = 'tree-checkbox';
        cb.checked = state.selectedFiles.has(node.path);
        cb.addEventListener('change', (e) => {
            e.stopPropagation();
            if (cb.checked) {
                state.selectedFiles.add(node.path);
            } else {
                state.selectedFiles.delete(node.path);
            }
            updateStats();
        });

        const icon = document.createElement('span');
        icon.className = 'tree-icon';
        icon.textContent = '\uD83D\uDCC4';

        const name = document.createElement('span');
        name.className = 'tree-name';
        name.textContent = node.name;

        item.appendChild(toggle);
        item.appendChild(cb);
        item.appendChild(icon);
        item.appendChild(name);

        item.addEventListener('click', (e) => {
            if (e.target === cb) return;
            cb.checked = !cb.checked;
            cb.dispatchEvent(new Event('change'));
        });

        div.appendChild(item);
    }

    return div;
}

function hasVisibleChildren(node) {
    if (!node.children) return false;
    return node.children.some(child => {
        if (!child.isDir) return child.path.toLowerCase().includes(state.searchFilter);
        return hasVisibleChildren(child);
    });
}

// ============================================================
// Results rendering
// ============================================================

function renderResults() {
    const container = $('#resultsContainer');
    container.innerHTML = '';

    const selectedArr = Array.from(state.selectedFiles).sort();
    const deps = state.dependencies || [];
    const includes = state.includes || [];

    if (selectedArr.length === 0 && deps.length === 0) {
        container.innerHTML = '<div class="empty-state"><div class="icon">&#128270;</div><p>Select files from the tree and click Analyze</p></div>';
        $('#depCount').textContent = '';
        return;
    }

    if (selectedArr.length > 0) {
        const section = document.createElement('div');
        section.className = 'section';
        section.innerHTML = `<div class="section-title">
            <span class="badge badge-blue">Selected</span>
            <span>${selectedArr.length} files</span>
        </div>`;

        selectedArr.forEach(path => {
            const item = document.createElement('div');
            item.className = 'file-item';
            item.innerHTML = `<span class="file-path">${escHtml(path)}</span>`;
            section.appendChild(item);
        });

        container.appendChild(section);
    }

    if (deps.length > 0) {
        const section = document.createElement('div');
        section.className = 'section';
        section.innerHTML = `<div class="section-title">
            <span class="badge badge-orange">Dependencies</span>
            <span>${deps.length} files</span>
        </div>`;

        deps.forEach(dep => {
            const item = document.createElement('div');
            item.className = 'file-item';
            item.innerHTML = `
                <span class="file-path">${escHtml(dep.filePath)}</span>
                <span class="file-ref">${escHtml(dep.className)} (${dep.refType}) from ${escHtml(shortPath(dep.referencedBy))}</span>
            `;
            section.appendChild(item);
        });

        container.appendChild(section);
    }

    if (includes.length > 0) {
        const section = document.createElement('div');
        section.className = 'section';
        section.innerHTML = `<div class="section-title">
            <span class="badge badge-gray">Include/Require</span>
            <span>${includes.length} references</span>
        </div>`;

        includes.forEach((inc, idx) => {
            const item = document.createElement('div');
            item.className = 'include-item';

            const cb = document.createElement('input');
            cb.type = 'checkbox';
            cb.checked = state.checkedIncludes.has(idx);
            cb.addEventListener('change', () => {
                if (cb.checked) state.checkedIncludes.add(idx);
                else state.checkedIncludes.delete(idx);
                updateCopyButton();
            });

            const typeSpan = document.createElement('span');
            typeSpan.className = 'badge badge-gray';
            typeSpan.textContent = inc.type;
            typeSpan.style.fontSize = '0.8em';

            const pathSpan = document.createElement('span');
            pathSpan.className = 'file-path';
            pathSpan.textContent = inc.resolved || inc.rawPath;

            const srcSpan = document.createElement('span');
            srcSpan.className = 'file-ref';
            srcSpan.textContent = `from ${shortPath(inc.sourceFile)}`;

            item.appendChild(cb);
            item.appendChild(typeSpan);
            item.appendChild(pathSpan);
            item.appendChild(srcSpan);

            section.appendChild(item);
        });

        container.appendChild(section);
    }

    const totalDeps = deps.length;
    $('#depCount').textContent = totalDeps > 0 ? `${totalDeps} deps` : '';
}

// ============================================================
// Helpers
// ============================================================

function getAllCopyFiles() {
    const files = new Set(state.selectedFiles);
    (state.dependencies || []).forEach(dep => files.add(dep.filePath));
    state.checkedIncludes.forEach(idx => {
        const inc = state.includes[idx];
        if (inc && inc.resolved) files.add(inc.resolved);
    });
    return Array.from(files);
}

function updateCopyButton() {
    $('#btnCopy').disabled = getAllCopyFiles().length === 0;
}

function updateStats() {
    const selected = state.selectedFiles.size;
    const deps = (state.dependencies || []).length;
    const total = getAllCopyFiles().length;
    $('#statusStats').textContent = `Selected: ${selected} | Dependencies: ${deps} | Total: ${total}`;
    updateCopyButton();
}

function setStatus(msg) {
    $('#statusText').textContent = msg;
}

function shortPath(p) {
    if (!p) return '';
    const parts = p.split('/');
    if (parts.length <= 2) return p;
    return '.../' + parts.slice(-2).join('/');
}

function escHtml(s) {
    if (!s) return '';
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

// ============================================================
// Resize handle
// ============================================================

(function initResize() {
    const handle = $('#resizeHandle');
    const panel = $('#panelLeft');
    let startX, startWidth;

    handle.addEventListener('mousedown', (e) => {
        startX = e.clientX;
        startWidth = panel.offsetWidth;
        handle.classList.add('active');
        document.addEventListener('mousemove', onMouseMove);
        document.addEventListener('mouseup', onMouseUp);
        e.preventDefault();
    });

    function onMouseMove(e) {
        const diff = e.clientX - startX;
        const newWidth = Math.max(200, Math.min(800, startWidth + diff));
        panel.style.width = newWidth + 'px';
    }

    function onMouseUp() {
        handle.classList.remove('active');
        document.removeEventListener('mousemove', onMouseMove);
        document.removeEventListener('mouseup', onMouseUp);
    }
})();

// ============================================================
// Init
// ============================================================

setStatus('Ready - select a project directory to begin');
