import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import TasksTable from '../../components/tasks/tasksTable';
import LoadingSpinner from '../../components/loadingSpinner/loadingSpinner';
import classes from "./mainPage.module.css";
import { Button } from "antd";

const MainPage = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        async function fetchData() {
            try {
                setLoading(true);
                setError(null);

                const response = await fetch(`https://todoappmoon.ru/v1/private/tasks`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: "include"
                });

                if (!response.ok) {
                    if (response.status === 403) {
                        navigate('/auth');
                        return;
                    }
                    throw new Error(`Ошибка загрузки: ${response.status}`);
                }

                const jsonData = await response.json();
                setData(jsonData.tasks || []);
            } catch (e) {
                setError(e.message);
                console.error('Ошибка загрузки задач:', e);
            } finally {
                setLoading(false);
            }
        }

        fetchData();
    }, [navigate]);

    const handleLogout = async () => {
        if (window.confirm('Вы уверены, что хотите выйти из профиля?')) {
            try {
                const response = await fetch('https://todoappmoon.ru/v1/auth/sign-out', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include'
                });

                if (response.ok) {
                    console.log('Выход выполнен успешно');
                    alert('Вы успешно вышли из системы');
                } else {
                    console.error('Ошибка при выходе:', response.status);
                    alert('Произошла ошибка при выходе из системы');
                }
            } catch (error) {
                console.error('Ошибка сети при выходе:', error);
                alert('Произошла ошибка сети при выходе из системы');
            } finally {
                navigate('/auth');
            }
        }
    };

    if (loading) {
        return (
            <div className={classes.loadingContainer}>
                <LoadingSpinner />
                <p>Uploading your tasks...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className={classes.errorContainer}>
                <div className={classes.errorContent}>
                    <h3>Cant upload your tasks</h3>
                    <p>{error}</p>
                    <Button
                        type="primary"
                        onClick={() => window.location.reload()}
                        className={classes.retryButton}
                    >
                        Try again!
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div className={classes.mainContainer}>
            <header className={classes.pageHeader}>
                <h1>My ToDoApp</h1>
                <div className={classes.headerButtons}>
                    <Button
                        type="primary"
                        onClick={() => navigate('/profile')}
                    >
                        Me
                    </Button>
                    <Button
                        type="default"
                        onClick={handleLogout}
                    >
                        Log Out
                    </Button>
                </div>
            </header>
            <main className={classes.pageContent}>
                <TasksTable tasks={data} />
            </main>
        </div>
    );
}

export default MainPage;