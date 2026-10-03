export const checkRequired = (fields) => {
    return fields.every((field) => {
        if (typeof field === 'string') return field.trim() !== '';
        if (typeof field === 'number') return !isNaN(field);
        return field !== null && field !== undefined;
    })
};
