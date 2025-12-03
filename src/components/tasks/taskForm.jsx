import { useState, useEffect } from 'react';
import { Form, Input, Button, DatePicker } from 'antd';
import dayjs from 'dayjs';
import classes from "./tasks.module.css";

const TaskForm = ({ onSubmit, onCancel, isEditing = false, task = null, parentId = null }) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (isEditing && task) {
            form.setFieldsValue({
                title: task.title,
                description: task.description,
                start_at: task.start_at ? dayjs(task.start_at) : null,
                end_at: task.end_at ? dayjs(task.end_at) : null,
            });
        } else {
            form.resetFields();
        }
    }, [isEditing, task, form]);

    const handleSubmit = async (values) => {
        setLoading(true);
        try {
            const url = isEditing 
                ? `https://todoappmoon.ru/v1/private/tasks/${task.id}`
                : 'https://todoappmoon.ru/v1/private/tasks';
                
            const method = isEditing ? 'PATCH' : 'POST';

            const requestBody = {
                title: values.title,
                description: values.description || "",
                start_at: values.start_at ? values.start_at.toISOString() : null,
                end_at: values.end_at ? values.end_at.toISOString() : null,
            };

            if (!isEditing && parentId) {
                requestBody.parent_id = parentId;
            }

            console.log('Sending data:', requestBody);

            const response = await fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify(requestBody)
            });

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Server error:', errorText);
                throw new Error('Network response was not ok');
            }

            const result = await response.json();
            onSubmit(result);
            form.resetFields();

        } catch (error) {
            console.error('Error saving task:', error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className={classes.TaskForm}>
            <h3>{isEditing ? 'Edit Task' : 'Create New Task'}</h3>
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
            >
                <Form.Item
                    label="Title"
                    name="title"
                    rules={[{ required: true, message: 'Please enter task title' }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="Description"
                    name="description"
                >
                    <Input.TextArea rows={4} />
                </Form.Item>

                <Form.Item
                    label="Start Date"
                    name="start_at"
                >
                    <DatePicker 
                        showTime 
                        format="YYYY-MM-DD HH:mm"
                        style={{ width: '100%' }}
                    />
                </Form.Item>

                <Form.Item
                    label="End Date"
                    name="end_at"
                >
                    <DatePicker 
                        showTime 
                        format="YYYY-MM-DD HH:mm"
                        style={{ width: '100%' }}
                    />
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading}>
                        {isEditing ? 'Update Task' : 'Create Task'}
                    </Button>
                    <Button onClick={onCancel} style={{ marginLeft: 8 }}>
                        Cancel
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
};

export default TaskForm;