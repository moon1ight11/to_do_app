import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, message } from "antd";
import classes from "../auth.module.css";

const RegForm = () => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const sendNewUser = async (userData) => {
        try {
            const response = await fetch('http://localhost:8080/v1/auth/sign-up', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(userData),
                credentials: "include"
            });

            if (response.ok) {
                message.success('Account created!');
                setError('');
                setTimeout(() => navigate('/'), 600);
                return await response.json();
            } else {
                const errorData = await response.json();
                const msg = errorData.error || 'Registration failed';
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
            await sendNewUser({
                user_name: values.userName,
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
            <h2>Create account</h2>
            <p>Start managing your tasks</p>

            {error && <div className={classes.errorMessage}>{error}</div>}

            <Form
                form={form}
                name="register"
                labelCol={{ span: 24 }}
                wrapperCol={{ span: 24 }}
                onFinish={onFinish}
                autoComplete="off"
                validateTrigger="onSubmit"
            >
                <Form.Item
                    label="Name"
                    name="userName"
                    validateTrigger="onSubmit"
                    rules={[
                        { required: true, message: 'Please enter your name' },
                        { min: 3, message: 'Name must be at least 3 characters' }
                    ]}
                >
                    <Input />
                </Form.Item>

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
                        { required: true, message: 'Please enter your password' },
                        { min: 6, message: 'Password must be at least 6 characters' }
                    ]}
                >
                    <Input.Password />
                </Form.Item>

                <Form.Item
                    label="Confirm Password"
                    name="confirmPassword"
                    dependencies={['userPass']}
                    validateTrigger="onSubmit"
                    rules={[
                        { required: true, message: 'Please confirm your password' },
                        ({ getFieldValue }) => ({
                            validator(_, value) {
                                if (!value || getFieldValue('userPass') === value) {
                                    return Promise.resolve();
                                }
                                return Promise.reject(new Error('Passwords do not match!'));
                            },
                        }),
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
                        Create Account
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
}

export default RegForm;