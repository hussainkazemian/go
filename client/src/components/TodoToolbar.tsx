import { Flex, Input, Select, Box } from "@chakra-ui/react";
import { useDeferredValue, useEffect, useState } from "react";

export type Filters = {
    status: "all" | "completed" | "active";
    sortBy: "createdAt" | "body" | "completed";
    order: "asc" | "desc";
    search: string;
};

export default function FilterBar({ value, onChange }: { value: Filters; onChange: (v: Filters) => void }) {
    const [local, setLocal] = useState<Filters>(value);
    const deferredSearch = useDeferredValue(local.search);

    useEffect(() => setLocal(value), [value]);

    useEffect(() => {
        onChange({ ...local, search: deferredSearch });
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [deferredSearch, local.status, local.sortBy, local.order]);


    return (
        <Box my={3}>
            <Flex gap={3} alignItems="center" justify="space-between" wrap="wrap">
                <Flex gap={2} alignItems="center" wrap="wrap">
                    <Select value={local.status} onChange={(e) => setLocal((s) => ({ ...s, status: e.target.value as Filters["status"] }))} maxW="180px">
                        <option value="all">All</option>
                        <option value="active">Active</option>
                        <option value="completed">Completed</option>
                    </Select>
                    <Select value={local.sortBy} onChange={(e) => setLocal((s) => ({ ...s, sortBy: e.target.value as Filters["sortBy"] }))} maxW="200px">
                        <option value="createdAt">Sort: Created</option>
                        <option value="body">Sort: Title</option>
                        <option value="completed">Sort: Status</option>
                    </Select>
                    <Select value={local.order} onChange={(e) => setLocal((s) => ({ ...s, order: e.target.value as Filters["order"] }))} maxW="180px">
                        <option value="desc">Order: Descending</option>
                        <option value="asc">Order: Ascending</option>
                    </Select>
                </Flex>
                <Input
                    placeholder="Search..."
                    value={local.search}
                    onChange={(e) => setLocal((s) => ({ ...s, search: e.target.value }))}
                    maxW={{ base: "100%", md: "320px" }}
                />
            </Flex>
        </Box>
    );
}
