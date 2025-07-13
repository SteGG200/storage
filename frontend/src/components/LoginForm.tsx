"use client";

import { Lock } from "lucide-react";
import { FormEvent, useState } from "react";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { login } from "@/lib/action/auth";
import { useQueryClient } from "@tanstack/react-query";

interface LoginFormProps {
  path: string;
}

export default function LoginForm({ path }: LoginFormProps) {
  const [errorMessage, setErrorMessage] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const queryClient = useQueryClient()

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
		setIsLoading(true);
		try {
			const formData = new FormData(event.currentTarget)
			await login(path, formData)
      await queryClient.invalidateQueries({
        queryKey: ['checkAuth', path]
      })
			setIsLoading(false)
		}catch(err){
			setIsLoading(false)
			setErrorMessage((err as Error).message)
		}
	};

  return (
    <form onSubmit={submit} className="space-y-4">
      <div className="space-y-2">
        <div className="relative">
          <Lock className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
          <Input
            type="password"
            name="password"
            placeholder="Enter password"
            value={password}
            onChange={(event) => {
              setPassword(event.target.value);
            }}
            className="pl-10 bg-gray-700 border-gray-600 text-gray-100 placeholder-gray-400"
            required
          />
        </div>
        {errorMessage.length > 0 && <p className="text-sm text-red-400">{errorMessage}</p>}
      </div>
      <Button
        type="submit"
        className="w-full bg-blue-600 hover:bg-blue-700 text-white"
        disabled={isLoading || !password.trim()}
      >
        {isLoading ? "Loging in..." : "Log in"}
      </Button>
    </form>
  );
}
