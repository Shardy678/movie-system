"use client"

import { useState } from "react"
import LoginForm from "./LoginForm"
import RegisterForm from "./RegisterForm"
import styles from "./LoginPage.module.css"

function LoginPage({ isLogin }: { isLogin: boolean }) {
  const [isLoginView, setIsLoginView] = useState(isLogin)

  const toggleView = () => {
    setIsLoginView(!isLoginView)
  }

  return (
    <div className={styles.authContainer}>
      <h1 className={styles.title}>{isLoginView ? "Login" : "Create an account"}</h1>
      {isLoginView ? <LoginForm /> : <RegisterForm />}
      <div className={styles.toggleContainer}>
        {isLoginView ? "Don't have an account?" : "Already have an account?"}
        <button onClick={toggleView} className={styles.toggleButton}>
          {isLoginView ? "Sign up" : "Log in"}
        </button>
      </div>
    </div>
  )
}

export default LoginPage

