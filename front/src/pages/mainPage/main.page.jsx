import { useState, useEffect, useRef, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import classes from "./mainPage.module.css";

const MainPage = () => {
    const [tasks, setTasks] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [showForm, setShowForm] = useState(false);
    const [editingTask, setEditingTask] = useState(null);
    const [editingSubtask, setEditingSubtask] = useState(null);
    const [form, setForm] = useState({ title: '', description: '', start_at: '', end_at: '' });
    const [subtaskInputs, setSubtaskInputs] = useState({});
    const [openMenu, setOpenMenu] = useState(null);
    const [openSubtaskMenu, setOpenSubtaskMenu] = useState(null);
    const menuRef = useRef(null);
    const navigate = useNavigate();
    const API = 'http://localhost:8080';

    useEffect(() => {
        const handler = (e) => {
            if (menuRef.current && !menuRef.current.contains(e.target)) {
                setOpenMenu(null);
                setOpenSubtaskMenu(null);
            }
        };
        document.addEventListener('mousedown', handler);
        return () => document.removeEventListener('mousedown', handler);
    }, []);

    const fetchTasks = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);
            const res = await fetch(`${API}/v1/private/tasks`, { credentials: 'include' });
            if (res.status === 403) return navigate('/auth');
            if (!res.ok) throw new Error('Failed to load tasks');
            const data = await res.json();
            setTasks(data.tasks || []);
        } catch (e) {
            setError(e.message);
        } finally {
            setLoading(false);
        }
    }, [API, navigate]);

    useEffect(() => { fetchTasks(); }, [fetchTasks]);

    const handleLogout = async () => {
        await fetch(`${API}/v1/private/sign-out`, { method: 'POST', credentials: 'include' });
        navigate('/auth');
    };

    const getDefaultDates = () => {
        const now = new Date();
        const later = new Date(now.getTime() + 60 * 60 * 1000);
        return {
            start: now.toISOString().slice(0, 16),
            end: later.toISOString().slice(0, 16)
        };
    };

    const openCreateForm = () => {
        const d = getDefaultDates();
        setEditingTask(null);
        setEditingSubtask(null);
        setForm({ title: '', description: '', start_at: d.start, end_at: d.end });
        setShowForm(true);
    };

    const openEditForm = (task) => {
        setEditingTask(task);
        setEditingSubtask(null);
        setForm({
            title: task.title || '',
            description: task.description || '',
            start_at: task.start_at ? task.start_at.slice(0, 16) : '',
            end_at: task.end_at ? task.end_at.slice(0, 16) : ''
        });
        setShowForm(true);
    };

    const openEditSubtask = (parentTask, subtask) => {
        setEditingTask(parentTask);
        setEditingSubtask(subtask);
        setForm({
            title: subtask.title || '',
            description: subtask.description || '',
            start_at: subtask.start_at ? subtask.start_at.slice(0, 16) : '',
            end_at: subtask.end_at ? subtask.end_at.slice(0, 16) : ''
        });
        setShowForm(true);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        const body = {
            title: form.title,
            description: form.description,
            start_at: form.start_at ? new Date(form.start_at).toISOString() : null,
            end_at: form.end_at ? new Date(form.end_at).toISOString() : null
        };

        let url, method;

        if (editingSubtask) {
            url = `${API}/v1/private/tasks/${editingSubtask.task_id}`;
            method = 'PATCH';
        } else if (editingTask) {
            url = `${API}/v1/private/tasks/${editingTask.task_id}`;
            method = 'PATCH';
        } else {
            url = `${API}/v1/private/tasks`;
            method = 'POST';
        }

        const res = await fetch(url, {
            method,
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(body)
        });
        if (res.ok) { setShowForm(false); setEditingSubtask(null); fetchTasks(); }
    };

    const handleDelete = async (taskId) => {
        if (!window.confirm('Delete this task?')) return;
        await fetch(`${API}/v1/private/tasks/${taskId}`, { method: 'DELETE', credentials: 'include' });
        setOpenMenu(null);
        fetchTasks();
    };

    const toggleComplete = async (task) => {
        const newStatus = !task.completed_at;
        await fetch(`${API}/v1/private/tasks/${task.task_id}`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ completed_at: newStatus })
        });
        setTasks(prev => prev.map(t =>
            t.task_id === task.task_id ? { ...t, completed_at: newStatus ? new Date().toISOString() : null } : t
        ));
    };

    const toggleSubtaskComplete = async (parentTaskId, subtask) => {
        const newStatus = !subtask.completed_at;
        await fetch(`${API}/v1/private/tasks/${subtask.task_id}`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ completed_at: newStatus })
        });
        setTasks(prev => prev.map(t => {
            if (t.task_id !== parentTaskId) return t;
            return {
                ...t,
                subtasks: (t.subtasks || []).map(s =>
                    s.task_id === subtask.task_id ? { ...s, completed_at: newStatus ? new Date().toISOString() : null } : s
                )
            };
        }));
    };

    const createSubtask = async (parentTaskId, title) => {
        if (!title.trim()) return;
        await fetch(`${API}/v1/private/tasks`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ title: title.trim(), parent_id: parentTaskId })
        });
        setSubtaskInputs(prev => ({ ...prev, [parentTaskId]: '' }));
        fetchTasks();
    };

    const deleteSubtask = async (subtaskId) => {
        await fetch(`${API}/v1/private/tasks/${subtaskId}`, { method: 'DELETE', credentials: 'include' });
        setOpenSubtaskMenu(null);
        fetchTasks();
    };

    const formatDate = (d) => {
        if (!d) return '';
        return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
    };

    const setQuickDate = (field, minutes) => {
        const d = new Date(Date.now() + minutes * 60 * 1000);
        setForm(prev => ({ ...prev, [field]: d.toISOString().slice(0, 16) }));
    };

    if (loading) return (
        <div className={classes.mainContainer}>
            <div className={classes.loadingContainer}>
                <div className={classes.dotsSpinner}>
                    <div className={classes.dot}></div>
                    <div className={classes.dot}></div>
                    <div className={classes.dot}></div>
                </div>
                <span className={classes.loadingText}>Loading tasks...</span>
            </div>
        </div>
    );

    if (error) return (
        <div className={classes.mainContainer}>
            <div className={classes.errorContainer}>
                <div className={classes.errorContent}>
                    <h3>Something went wrong</h3>
                    <p>{error}</p>
                    <button className={classes.retryButton} onClick={fetchTasks}>Try again</button>
                </div>
            </div>
        </div>
    );

    return (
        <div className={classes.mainContainer}>
            <div className={classes.pageHeader}>
                <h1>My Tasks</h1>
                <span className={classes.taskCount}>{tasks.length} {tasks.length === 1 ? 'task' : 'tasks'}</span>
            </div>

            {tasks.length === 0 ? (
                <div className={classes.emptyState}>
                    <div className={classes.emptyIcon}>📋</div>
                    <h3 className={classes.emptyTitle}>No tasks yet</h3>
                    <p className={classes.emptySubtitle}>Create your first task to get started</p>
                    <button className={classes.emptyBtn} onClick={openCreateForm}>Create Task</button>
                </div>
            ) : (
                <div className={classes.taskList}>
                    {tasks.map(task => (
                        <div key={task.task_id}
                            className={`${classes.taskCard} ${task.completed_at ? classes.completed : classes.pending}`}>

                            <div className={classes.taskHeader}>
                                <div className={`${classes.checkbox} ${task.completed_at ? classes.checked : ''}`}
                                    onClick={() => toggleComplete(task)} />
                                <h3 className={classes.taskTitle}>{task.title}</h3>
                                <div style={{ position: 'relative' }}>
                                    <button className={classes.menuBtn}
                                        onClick={() => setOpenMenu(openMenu === task.task_id ? null : task.task_id)}>
                                        <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                                            <circle cx="5" cy="12" r="2" /><circle cx="12" cy="12" r="2" /><circle cx="19" cy="12" r="2" />
                                        </svg>
                                    </button>
                                    {openMenu === task.task_id && (
                                        <div className={classes.dropdown} ref={menuRef}>
                                            <button className={classes.dropdownItem}
                                                onClick={() => { setOpenMenu(null); openEditForm(task); }}>
                                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" /><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" /></svg>
                                                Edit
                                            </button>
                                            <button className={`${classes.dropdownItem} ${classes.danger}`}
                                                onClick={() => handleDelete(task.task_id)}>
                                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /></svg>
                                                Delete
                                            </button>
                                        </div>
                                    )}
                                </div>
                            </div>

                            {task.description && <p className={classes.taskDescription}>{task.description}</p>}

                            <div className={classes.taskMeta}>
                                {task.start_at && <span className={classes.taskTime}>Start: {formatDate(task.start_at)}</span>}
                                {task.end_at && <span className={classes.taskTime}>End: {formatDate(task.end_at)}</span>}
                            </div>

                            {task.subtasks && task.subtasks.length > 0 && (
                                <div className={classes.subtasksSection}>
                                    {task.subtasks.map(sub => (
                                        <div key={sub.task_id} className={classes.subtaskItem}>
                                            <div className={`${classes.checkbox} ${sub.completed_at ? classes.checked : ''}`}
                                                onClick={() => toggleSubtaskComplete(task.task_id, sub)} />
                                            <span className={`${classes.subtaskTitle} ${sub.completed_at ? classes.done : ''}`}>{sub.title}</span>
                                            <div style={{ position: 'relative' }}>
                                                <button className={classes.subtaskMenuBtn}
                                                    onClick={() => setOpenSubtaskMenu(openSubtaskMenu === sub.task_id ? null : sub.task_id)}>
                                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                                                        <circle cx="5" cy="12" r="1.5" /><circle cx="12" cy="12" r="1.5" /><circle cx="19" cy="12" r="1.5" />
                                                    </svg>
                                                </button>
                                                {openSubtaskMenu === sub.task_id && (
                                                    <div className={classes.subtaskDropdown}>
                                                        <button className={classes.dropdownItem}
                                                            onClick={() => { setOpenSubtaskMenu(null); openEditSubtask(task, sub); }}>
                                                            Edit
                                                        </button>
                                                        <button className={`${classes.dropdownItem} ${classes.danger}`}
                                                            onClick={() => deleteSubtask(sub.task_id)}>
                                                            Delete
                                                        </button>
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}

                            <div className={classes.addSubtaskRow} style={{ marginLeft: 36 }}>
                                <input
                                    placeholder="Add subtask..."
                                    value={subtaskInputs[task.task_id] || ''}
                                    onChange={e => setSubtaskInputs(prev => ({ ...prev, [task.task_id]: e.target.value }))}
                                    onKeyDown={e => {
                                        if (e.key === 'Enter') {
                                            e.preventDefault();
                                            createSubtask(task.task_id, subtaskInputs[task.task_id] || '');
                                        }
                                    }}
                                />
                                <button onClick={() => createSubtask(task.task_id, subtaskInputs[task.task_id] || '')}>
                                    + Add
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {tasks.length > 0 && (
                <button className={classes.fab} onClick={openCreateForm}>
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                        <line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" />
                    </svg>
                </button>
            )}

            <div className={classes.bottomNav}>
                <button className={`${classes.navBtn} ${classes.active}`} title="Tasks">
                    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <rect x="3" y="3" width="7" height="7" rx="1" /><rect x="14" y="3" width="7" height="7" rx="1" />
                        <rect x="3" y="14" width="7" height="7" rx="1" /><rect x="14" y="14" width="7" height="7" rx="1" />
                    </svg>
                </button>
                <button className={classes.navBtn} onClick={() => navigate('/profile')} title="Profile">
                    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" /><circle cx="12" cy="7" r="4" />
                    </svg>
                </button>
                <button className={classes.navBtn} onClick={handleLogout} title="Logout">
                    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" /><polyline points="16 17 21 12 16 7" /><line x1="21" y1="12" x2="9" y2="12" />
                    </svg>
                </button>
            </div>

            {showForm && (
                <>
                    <div className={classes.overlay} onClick={() => { setShowForm(false); setEditingSubtask(null); }}></div>
                    <div className={classes.bottomSheet}>
                        <div className={classes.sheetHandle}></div>
                        <h2 className={classes.sheetTitle}>
                            {editingSubtask ? 'Edit Subtask' : editingTask ? 'Edit Task' : 'New Task'}
                        </h2>
                        <form className={classes.sheetForm} onSubmit={handleSubmit}>
                            <div className={classes.formGroup}>
                                <label className={classes.formLabel}>Title</label>
                                <input className={classes.formInput} value={form.title} required
                                    onChange={e => setForm(prev => ({ ...prev, title: e.target.value }))}
                                    placeholder="What needs to be done?" />
                            </div>
                            <div className={classes.formGroup}>
                                <label className={classes.formLabel}>Description</label>
                                <textarea className={classes.formTextarea} rows={2} value={form.description}
                                    onChange={e => setForm(prev => ({ ...prev, description: e.target.value }))}
                                    placeholder="Add details..." />
                            </div>
                            <div className={classes.formGroup}>
                                <label className={classes.formLabel}>Start</label>
                                <input className={classes.formInput} type="datetime-local" value={form.start_at}
                                    onChange={e => setForm(prev => ({ ...prev, start_at: e.target.value }))} />
                                <div className={classes.dateQuick}>
                                    <span className={classes.dateChip} onClick={() => setQuickDate('start_at', 0)}>Now</span>
                                    <span className={classes.dateChip} onClick={() => setQuickDate('start_at', 30)}>+30m</span>
                                    <span className={classes.dateChip} onClick={() => setQuickDate('start_at', 60)}>+1h</span>
                                    <span className={classes.dateChip} onClick={() => setQuickDate('start_at', 1440)}>+24h</span>
                                </div>
                            </div>
                            <div className={classes.formGroup}>
                                <label className={classes.formLabel}>End</label>
                                <input className={classes.formInput} type="datetime-local" value={form.end_at}
                                    onChange={e => setForm(prev => ({ ...prev, end_at: e.target.value }))} />
                                <div className={classes.dateQuick}>
                                    <span className={classes.dateChip} onClick={() => setQuickDate('end_at', 30)}>+30m</span>
                                    <span className={classes.dateChip} onClick={() => setQuickDate('end_at', 60)}>+1h</span>
                                    <span className={classes.dateChip} onClick={() => setQuickDate('end_at', 120)}>+2h</span>
                                    <span className={classes.dateChip} onClick={() => setQuickDate('end_at', 1440)}>+24h</span>
                                </div>
                            </div>
                            <div className={classes.formActions}>
                                <button type="submit" className={`${classes.btn} ${classes.btnPrimary}`}>
                                    {editingTask || editingSubtask ? 'Save Changes' : 'Create Task'}
                                </button>
                                <button type="button" className={`${classes.btn} ${classes.btnCancel}`}
                                    onClick={() => { setShowForm(false); setEditingSubtask(null); }}>Cancel</button>
                            </div>
                        </form>
                    </div>
                </>
            )}
        </div>
    );
};

export default MainPage;