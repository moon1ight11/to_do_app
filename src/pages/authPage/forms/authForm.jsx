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
            const response = await fetch('https://todoappmoon.ru/v1/auth/sign-in', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData),
                credentials: "include"
            });

            if (response.ok) {
                const result = await response.json();
                console.log('Успешная аутентификация:', result);
                message.success('Аутентификация прошла успешно!');
                setError('');

                setTimeout(() => {
                    navigate('/');
                }, 1000);

                return result;
            } else {
                const errorData = await response.json();
                console.error('Ошибка сервера:', response.status, errorData);
                
                const errorMessage = errorData.error || response.statusText;
                setError(errorMessage);
                message.error(`Ошибка аутентификации: ${errorMessage}`);
                
                throw new Error(errorMessage);
            }
        } catch (error) {
            console.error('Ошибка:', error);
            const errorMsg = error.message || 'Произошла ошибка при аутентификации';
            setError(errorMsg);
            message.error(errorMsg);
            throw error;
        }
    };

    const onFinish = async (values) => {
        setLoading(true);
        setError('');
        try {
            const User = {
                user_pass: values.userPass,
                user_email: values.userEmail
            };

            await sendUser(User);
            form.resetFields();

        } catch (error) {
        } finally {
            setLoading(false);
        }
    };

    const onFinishFailed = (errorInfo) => {
        console.log('Failed:', errorInfo);
        message.warning('Пожалуйста, заполните все поля правильно');
    };

    return (
        <div className={classes.formContainer}>
            {error && (
                <div className={classes.errorMessage} style={{
                    color: 'red',
                    textAlign: 'center',
                    marginBottom: '16px',
                    padding: '8px',
                    backgroundColor: '#fff2f0',
                    border: '1px solid #ffccc7',
                    borderRadius: '4px'
                }}>
                    {error}
                </div>
            )}
            
            <Form
                form={form}
                name="auth"
                labelCol={{ span: 10 }}
                wrapperCol={{ span: 17 }}
                style={{ maxWidth: 600 }}
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
                autoComplete="off"
                validateTrigger="onSubmit"
            >
                <Form.Item
                    label="Email"
                    name="userEmail"
                    validateTrigger="onSubmit"
                    rules={[
                        {
                            required: true,
                            message: 'Please, enter your email!'
                        },
                        {
                            type: 'email',
                            message: 'Please enter a valid email address!'
                        }
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="Password"
                    name="userPass"
                    validateTrigger="onSubmit"
                    rules={[
                        {
                            required: true,
                            message: 'Please, enter your password!'
                        },
                    ]}
                >
                    <Input.Password />
                </Form.Item>

                <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
                    <Button
                        type="primary"
                        htmlType="submit"
                        loading={loading}
                        className={classes.submitButton}
                    >
                        Sign IN!
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
}

export default AuthForm;