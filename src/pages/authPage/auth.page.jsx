import { useState } from 'react';
import classes from "./auth.module.css";
import AuthForm from './forms/authForm';
import RegForm from './forms/regForm';
import { Button } from "antd";

const AuthPage = () => {
    const [activeForm, setActiveForm] = useState('auth');

    return (
        <div className={classes.RegAuthPage}>
            <div className={classes.Sets}>
                <Button
                    type={`${activeForm === 'auth' ? "primary" : "default"}`}
                    onClick={() => setActiveForm('auth')}
                >
                    Authorization
                </Button>
                <Button
                    type={`${activeForm === 'reg' ? "primary" : "default"}`}
                    onClick={() => setActiveForm('reg')}
                >
                  Registration
                </Button>
            </div>
            <div>
                {activeForm === 'auth' ? <AuthForm /> : <RegForm />}
            </div>
        </div>
    );
}

export default AuthPage;
