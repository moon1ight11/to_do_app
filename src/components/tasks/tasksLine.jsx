import { ClockCircleOutlined, CalendarOutlined, DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import classes from "./tasks.module.css";
import { Button, message } from 'antd';
import { useState } from 'react';
import TaskForm from './taskForm';

const TasksLine = ({ title, description, start_at, end_at, completed_at, task_id, onDelete, onUpdate, subtasks }) => {
    const [isComplete, setIsComplete] = useState(!!completed_at);
    const [isEditing, setIsEditing] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [showSubtasks, setShowSubtasks] = useState(true);
    const [isCreatingSubtask, setIsCreatingSubtask] = useState(false);
    const [editingSubtask, setEditingSubtask] = useState(null);
    const [currentSubtasks, setCurrentSubtasks] = useState(subtasks || []);

    const formatDate = (dateString) => {
        if (!dateString) return 'Not set';
        try {
            const date = new Date(dateString);
            return new Intl.DateTimeFormat('en-US', {
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            }).format(date);
        } catch (error) {
            return 'Invalid date';
        }
    };

    const handleStatusToggle = async () => {
        setIsLoading(true);
        try {
            const newStatus = !isComplete;
            const response = await fetch(`http://localhost:8080/v1/private/tasks/${task_id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "include",
                body: JSON.stringify({
                    completed_at: newStatus
                })
            });

            if (response.ok) {
                setIsComplete(newStatus);
                message.success(`Task marked as ${newStatus ? 'completed' : 'in progress'}`);
            } else {
                throw new Error('Failed to update task status');
            }
        } catch (error) {
            console.error('Error updating task status:', error);
            message.error('Failed to update task status');
        } finally {
            setIsLoading(false);
        }
    };

    const handleDelete = async () => {
        try {
            const response = await fetch(`http://localhost:8080/v1/private/tasks/${task_id}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "include"
            });

            if (response.ok) {
                message.success('Task deleted successfully');
                if (onDelete) {
                    onDelete(task_id);
                }
            } else {
                throw new Error('Failed to delete task');
            }
        } catch (error) {
            console.error('Error deleting task:', error);
            message.error('Failed to delete task');
        }
    };

    const handleEdit = async (formData) => {
        try {
            const response = await fetch(`http://localhost:8080/v1/private/tasks/${task_id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "include",
                body: JSON.stringify(formData)
            });

            if (response.ok) {
                message.success('Task updated successfully');
                setIsEditing(false);
                if (onUpdate) {
                    onUpdate();
                }
            } else {
                throw new Error('Failed to update task');
            }
        } catch (error) {
            console.error('Error updating task:', error);
            message.error('Failed to update task');
        }
    };

    const handleCancelEdit = () => {
        setIsEditing(false);
    };

    const confirmDelete = () => {
        if (window.confirm('Are you sure you want to delete this task?')) {
            handleDelete();
        }
    };

    // Функции для работы с подзадачами
    const handleCreateSubtask = async (formData) => {
        try {
            message.success('Subtask created successfully');
            setIsCreatingSubtask(false);

            if (onUpdate) {
                onUpdate();
            }
        } catch (error) {
            console.error('Error creating subtask:', error);
            message.error('Failed to create subtask');
        }
    };

    const handleUpdateSubtask = async (formData) => {
        try {
            const response = await fetch(`http://localhost:8080/v1/private/tasks/${editingSubtask.task_id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "include",
                body: JSON.stringify(formData)
            });

            if (response.ok) {
                const updatedSubtask = await response.json();
                setCurrentSubtasks(prev =>
                    prev.map(subtask =>
                        subtask.task_id === editingSubtask.task_id
                            ? {
                                ...subtask,
                                ...updatedSubtask,
                                task_id: subtask.task_id
                            }
                            : subtask
                    )
                );
                message.success('Subtask updated successfully');
                setEditingSubtask(null);
                if (onUpdate) {
                    onUpdate();
                }

            } else {
                throw new Error('Failed to update subtask');
            }
        } catch (error) {
            console.error('Error updating subtask:', error);
            message.error('Failed to update subtask');
        }
    };

    const handleDeleteSubtask = async (subtaskId) => {
        try {
            const response = await fetch(`http://localhost:8080/v1/private/tasks/${subtaskId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "include"
            });

            if (response.ok) {
                setCurrentSubtasks(prev => prev.filter(subtask => subtask.task_id !== subtaskId));
                message.success('Subtask deleted successfully');
            } else {
                throw new Error('Failed to delete subtask');
            }
        } catch (error) {
            console.error('Error deleting subtask:', error);
            message.error('Failed to delete subtask');
        }
    };

    const toggleSubtaskStatus = async (subtaskId, currentCompletedAt) => {
        try {
            const newStatus = !currentCompletedAt;
            const response = await fetch(`http://localhost:8080/v1/private/tasks/${subtaskId}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "include",
                body: JSON.stringify({
                    completed_at: newStatus
                })
            });

            if (response.ok) {
                setCurrentSubtasks(prev =>
                    prev.map(subtask =>
                        subtask.task_id === subtaskId
                            ? { ...subtask, completed_at: newStatus }
                            : subtask
                    )
                );
                message.success(`Subtask marked as ${newStatus ? 'completed' : 'in progress'}`);
            } else {
                throw new Error('Failed to update subtask status');
            }
        } catch (error) {
            console.error('Error updating subtask status:', error);
            message.error('Failed to update subtask status');
        }
    };

    const confirmDeleteSubtask = (subtaskId, subtaskTitle) => {
        if (window.confirm(`Are you sure you want to delete subtask "${subtaskTitle}"?`)) {
            handleDeleteSubtask(subtaskId);
        }
    };

    const handleEditSubtask = (subtask) => {
        setEditingSubtask(subtask);
    };

    const handleCancelSubtaskEdit = () => {
        setEditingSubtask(null);
    };

    return (
        <>
            <div className={classes.TasksLine}>
                <div className={classes.TasksHeader}>
                    <h3 className={classes.TasksTitle}>
                        {title}
                    </h3>
                    <div className={classes.HeaderActions}>
                        <Button
                            type="text"
                            icon={<EditOutlined />}
                            onClick={() => setIsEditing(true)}
                            className={classes.EditButton}
                            title="Edit task"
                        />
                        <Button
                            type="text"
                            danger
                            icon={<DeleteOutlined />}
                            onClick={confirmDelete}
                            className={classes.DeleteButton}
                            title="Delete task"
                        />
                        <div
                            className={`${classes.StatusBadge} ${isComplete ? classes.complete : classes.incomplete} ${isLoading ? classes.loading : ''}`}
                            onClick={handleStatusToggle}
                            style={{ cursor: isLoading ? 'not-allowed' : 'pointer' }}
                            title={isLoading ? 'Updating...' : `Click to mark as ${isComplete ? 'in progress' : 'completed'}`}
                        >
                            {isLoading ? 'Updating...' : (isComplete ? 'Completed' : 'In Progress')}
                        </div>
                    </div>
                </div>

                {description && (
                    <div className={classes.TasksDescription}>
                        {description}
                    </div>
                )}

                <div className={classes.TasksMeta}>
                    <div className={classes.TimeSlot}>
                        <CalendarOutlined className={classes.TimeIcon} />
                        <span>Start: {formatDate(start_at)}</span>
                    </div>

                    <div className={classes.TimeSlot}>
                        <ClockCircleOutlined className={classes.TimeIcon} />
                        <span>End: {formatDate(end_at)}</span>
                    </div>
                </div>

                {/* Секция подзадач */}
                <div className={classes.SubtasksSection}>
                    <div className={classes.SubtasksHeader}>
                        <Button
                            type="text"
                            icon={<PlusOutlined />}
                            onClick={() => setIsCreatingSubtask(true)}
                            className={classes.AddSubtaskButton}
                        >
                            Add Subtask
                        </Button>

                        {currentSubtasks.length > 0 && (
                            <Button
                                type="text"
                                onClick={() => setShowSubtasks(!showSubtasks)}
                                className={classes.ToggleSubtasksButton}
                            >
                                {showSubtasks ? 'Hide' : 'Show'} Subtasks ({currentSubtasks.length})
                            </Button>
                        )}
                    </div>

                    {isCreatingSubtask && (
                        <div className={classes.SubtaskForm}>
                            <TaskForm
                                onSubmit={handleCreateSubtask}
                                onCancel={() => setIsCreatingSubtask(false)}
                                isEditing={false}
                                task={null}
                                parentId={task_id}
                            />
                        </div>
                    )}

                    {showSubtasks && currentSubtasks.length > 0 && (
                        <div className={classes.SubtasksList}>
                            {currentSubtasks.map((subtask) => (
                                <div key={subtask.task_id} className={classes.SubtaskItem}>
                                    {editingSubtask && editingSubtask.task_id === subtask.task_id ? (
                                        <div className={classes.SubtaskForm}>
                                            <TaskForm
                                                task={{
                                                    id: subtask.task_id,
                                                    title: subtask.title,
                                                    description: subtask.description,
                                                    start_at: subtask.start_at,
                                                    end_at: subtask.end_at,
                                                    completed_at: subtask.completed_at
                                                }}
                                                onSubmit={handleUpdateSubtask}
                                                onCancel={handleCancelSubtaskEdit}
                                                isEditing={true}
                                            />
                                        </div>
                                    ) : (
                                        <>
                                            <div className={classes.SubtaskHeader}>
                                                <div className={classes.SubtaskTitle}>
                                                    <span
                                                        className={`${classes.SubtaskStatus} ${subtask.completed_at ? classes.complete : classes.incomplete}`}
                                                        onClick={() => toggleSubtaskStatus(subtask.task_id, subtask.completed_at)}
                                                        title={`Click to mark as ${subtask.completed_at ? 'in progress' : 'completed'}`}
                                                    >
                                                        {subtask.completed_at ? '✓' : '○'}
                                                    </span>
                                                    <span className={subtask.completed_at ? classes.CompletedText : ''}>
                                                        {subtask.title}
                                                    </span>
                                                </div>
                                                <div className={classes.SubtaskActions}>
                                                    <Button
                                                        type="text"
                                                        size="small"
                                                        icon={<EditOutlined />}
                                                        onClick={() => handleEditSubtask(subtask)}
                                                        title="Edit subtask"
                                                    />
                                                    <Button
                                                        type="text"
                                                        size="small"
                                                        danger
                                                        icon={<DeleteOutlined />}
                                                        onClick={() => confirmDeleteSubtask(subtask.task_id, subtask.title)}
                                                        title="Delete subtask"
                                                    />
                                                </div>
                                            </div>

                                            {subtask.description && (
                                                <div className={classes.SubtaskDescription}>
                                                    {subtask.description}
                                                </div>
                                            )}

                                            <div className={classes.SubtaskMeta}>
                                                <div className={classes.TimeSlot}>
                                                    <CalendarOutlined className={classes.TimeIcon} />
                                                    <span>Start: {formatDate(subtask.start_at)}</span>
                                                </div>
                                                <div className={classes.TimeSlot}>
                                                    <ClockCircleOutlined className={classes.TimeIcon} />
                                                    <span>End: {formatDate(subtask.end_at)}</span>
                                                </div>
                                            </div>
                                        </>
                                    )}
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>

            {isEditing && (
                <TaskForm
                    task={{
                        id: task_id,
                        title,
                        description,
                        start_at,
                        end_at,
                        completed_at
                    }}
                    onSubmit={handleEdit}
                    onCancel={handleCancelEdit}
                    isEditing={true}
                />
            )}
        </>
    );
};

export default TasksLine;