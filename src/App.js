import { Routes, Route } from "react-router-dom";
import MainPage from "./pages/mainPage/main.page";
import AuthPage from "./pages/authPage/auth.page";
import ProfilePage from "./pages/profilePage/profilePage";

function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<MainPage />} />
        <Route path="/auth" element={<AuthPage />} />
        <Route path="/profile" element={<ProfilePage />} />
      </Routes>
    </>
  );
}

export default App;
