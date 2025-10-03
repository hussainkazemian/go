import { Box, Divider, Flex, Spinner, Stack, Text } from "@chakra-ui/react";

import TodoItem from "./TodoItem";
import { useQuery } from "@tanstack/react-query";
import { BASE_URL } from "../App";
import { useState, useMemo } from "react";
import FilterBar, { Filters } from "./TodoToolbar";
import TodoForm from "./TodoForm";

export type Todo = {
	_id: string; // Mongo ObjectID as string when serialized
	body: string;
	completed: boolean;
	createdAt: string;
};

const TodoList = () => {
	const [filters, setFilters] = useState<Filters>({ status: "all", sortBy: "createdAt", order: "desc", search: "" });

	const queryString = useMemo(() => {
		const params = new URLSearchParams();
		if (filters.status && filters.status !== "all") params.set("status", filters.status);
		if (filters.sortBy) params.set("sortBy", filters.sortBy);
		if (filters.order) params.set("order", filters.order);
		if (filters.search) params.set("search", filters.search);
		const qs = params.toString();
		return qs ? `?${qs}` : "";
	}, [filters]);

	const { data: todos, isLoading } = useQuery<Todo[]>({
		queryKey: ["todos", filters],
		queryFn: async () => {
			const res = await fetch(`${BASE_URL}/todos${queryString}`);
			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || "Something went wrong");
			}
			return data || [];
		},
	});

	return (
		<>
			
			<FilterBar value={filters} onChange={setFilters} />
			<Divider my={3} opacity={0.25} />
			<Text
				fontSize={"4xl"}
				textTransform={"uppercase"}
				fontWeight={"bold"}
				textAlign={"center"}
				my={2}
				bgGradient='linear(to-l, #0b85f8, #00ffff)'
				bgClip='text'
			>
				DAILY NOTES
			</Text>
			<Box my={2}>
				<TodoForm />
			</Box>
			{isLoading && (
				<Flex justifyContent={"center"} my={4}>
					<Spinner size={"xl"} />
				</Flex>
			)}

			{!isLoading && todos?.length === 0 && (
				<Stack alignItems={"center"} gap='3'>
					<Text fontSize={"xl"} textAlign={"center"} color={"gray.500"}>
						All tasks completed! 🤞
					</Text>
					<img src='/go.png' alt='Go logo' width={70} height={70} />
				</Stack>
			)}
			<Stack gap={3} mt={3}>
				{todos?.map((todo) => (
					<TodoItem key={todo._id} todo={todo} />
				))}
			</Stack>
		</>
	);
};
export default TodoList;

// STARTER CODE:

// import { Flex, Spinner, Stack, Text } from "@chakra-ui/react";
// import { useState } from "react";
// import TodoItem from "./TodoItem";

// const TodoList = () => {
// 	const [isLoading, setIsLoading] = useState(true);
// 	const todos = [
// 		{
// 			_id: 1,
// 			body: "Buy groceries",
// 			completed: true,
// 		},
// 		{
// 			_id: 2,
// 			body: "Walk the dog",
// 			completed: false,
// 		},
// 		{
// 			_id: 3,
// 			body: "Do laundry",
// 			completed: false,
// 		},
// 		{
// 			_id: 4,
// 			body: "Cook dinner",
// 			completed: true,
// 		},
// 	];
// 	return (
// 		<>
// 			<Text fontSize={"4xl"} textTransform={"uppercase"} fontWeight={"bold"} textAlign={"center"} my={2}>
// 				Today's Tasks
// 			</Text>
// 			{isLoading && (
// 				<Flex justifyContent={"center"} my={4}>
// 					<Spinner size={"xl"} />
// 				</Flex>
// 			)}
// 			{!isLoading && todos?.length === 0 && (
// 				<Stack alignItems={"center"} gap='3'>
// 					<Text fontSize={"xl"} textAlign={"center"} color={"gray.500"}>
// 						All tasks completed! 🤞
// 					</Text>
// 					<img src='/go.png' alt='Go logo' width={70} height={70} />
// 				</Stack>
// 			)}
// 			<Stack gap={3}>
// 				{todos?.map((todo) => (
// 					<TodoItem key={todo._id} todo={todo} />
// 				))}
// 			</Stack>
// 		</>
// 	);
// };
// export default TodoList;
