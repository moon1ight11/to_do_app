import { useState } from 'react';
import classes from "./auth.module.css";
import AuthForm from './forms/authForm';
import RegForm from './forms/regForm';

const AuthPage = () => {
    const [activeForm, setActiveForm] = useState('auth');

    return (
        <div className={classes.RegAuthPage}>
            <div className={`${classes.Sets} ${activeForm === 'reg' ? classes.right : ''}`}>
                <button
                    className={activeForm === 'auth' ? classes.tabActive : classes.tab}
                    onClick={() => setActiveForm('auth')}
                >
                    Sign In
                </button>
                <button
                    className={activeForm === 'reg' ? classes.tabActive : classes.tab}
                    onClick={() => setActiveForm('reg')}
                >
                    Sign Up
                </button>
            </div>
            <div>
                {activeForm === 'auth' ? <AuthForm /> : <RegForm />}
            </div>
        </div>
    );
}

export default AuthPage;