import classes from "./LoadingSpinner.module.css"

const LoadingSpinner = () => {
  return (
    <div className={classes.spinner}>
      <div className={classes.dot}></div>
      <div className={classes.dot}></div>
      <div className={classes.dot}></div>
    </div>
  );
};

export default LoadingSpinner;