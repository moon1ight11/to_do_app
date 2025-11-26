import classes from "./LoadingSpinner.module.css"

const LoadingSpinner = () => {
  return (
    <div className={classes.cubeSpinner}>
      <div className={classes.cube}>
        <div className={classes.front}></div>
        <div className={classes.back}></div>
        <div className={classes.left}></div>
        <div className={classes.right}></div>
        <div className={classes.top}></div>
        <div className={classes.bottom}></div>
      </div>
    </div>
  );
};

export default LoadingSpinner;
