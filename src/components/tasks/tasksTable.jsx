import React, { useState } from 'react';
import TasksLine from './tasksLine';
import LoadingSpinner from '../loadingSpinner/loadingSpinner';
import classes from "./tasks.module.css";
import { Button } from 'antd';
import TaskForm from './taskForm';

const TasksTable = ({ tasks, onTasksUpdate }) => {
    const [isFormVisible, setIsFormVisible] = useState(false);

    const handleCreateTask = () => {
        setIsFormVisible(true);
    };

    const handleCloseForm = () => {
        setIsFormVisible(false);
    };

    const handleFormSubmit = (taskData) => {
        console.log('Создаем задачу:', taskData);
        handleCloseForm();
        window.location.reload();
    };

    const handleDeleteTask = (deletedTaskId) => {
        if (onTasksUpdate) {
            onTasksUpdate();
        } else {
            window.location.reload();
        }
    };

    const handleTaskUpdate = () => {
        if (onTasksUpdate) {
            onTasksUpdate();
        } else {
            window.location.reload();
        }
    };

    return (
        <div className={classes.TasksTable}>
            <Button
                type="default"
                onClick={handleCreateTask}
            >
                Create New Task
            </Button>

            {isFormVisible && (
                <TaskForm
                    onSubmit={handleFormSubmit}
                    onCancel={handleCloseForm}
                />
            )}

            {tasks && tasks.length ? (
                tasks.map(task => (
                    <TasksLine
                        key={task.id || task.title}
                        title={task.title}
                        description={task.description}
                        start_at={task.start_at}
                        end_at={task.end_at}
                        completed_at={task.completed_at}
                        task_id={task.task_id}
                        subtasks={task.subtasks || []}
                        onDelete={handleDeleteTask}
                        onUpdate={handleTaskUpdate}
                    />
                ))
            ) : (
                <div className={classes.MYFT}>
                    <LoadingSpinner />
                    Create your first task
                </div>
            )}
        </div>
    );
}

export default TasksTable;