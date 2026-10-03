import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, message } from "antd";
import classes from "../auth.module.css";

const AuthForm = () => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const sendUser = async (userData) => {
        try {
            const response = await fetch('http://localhost:8080/v1/auth/sign-in', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(userData),
                credentials: "include"
            });

            if (response.ok) {
                message.success('Welcome back!');
                setError('');
                setTimeout(() => navigate('/'), 600);
                return await response.json();
            } else {
                const errorData = await response.json();
                const msg = errorData.error || 'Invalid credentials';
                setError(msg);
                message.error(msg);
                throw new Error(msg);
            }
        } catch (error) {
            const msg = error.message || 'Network error';
            setError(msg);
            message.error(msg);
            throw error;
        }
    };

    const onFinish = async (values) => {
        setLoading(true);
        setError('');
        try {
            await sendUser({
                user_pass: values.userPass,
                user_email: values.userEmail
            });
            form.resetFields();
        } catch (error) {
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className={classes.formContainer}>
            <h2>Welcome back</h2>
            <p>Sign in to your account</p>

            {error && <div className={classes.errorMessage}>{error}</div>}

            <Form
                form={form}
                name="auth"
                labelCol={{ span: 24 }}
                wrapperCol={{ span: 24 }}
                onFinish={onFinish}
                autoComplete="off"
                validateTrigger="onSubmit"
            >
                <Form.Item
                    label="Email"
                    name="userEmail"
                    validateTrigger="onSubmit"
                    rules={[
                        { required: true, message: 'Please enter your email' },
                        { type: 'email', message: 'Please enter a valid email' }
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="Password"
                    name="userPass"
                    validateTrigger="onSubmit"
                    rules={[
                        { required: true, message: 'Please enter your password' }
                    ]}
                >
                    <Input.Password />
                </Form.Item>

                <Form.Item>
                    <Button
                        type="primary"
                        htmlType="submit"
                        loading={loading}
                        className={classes.submitButton}
                    >
                        Sign In
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
}

export default AuthForm;