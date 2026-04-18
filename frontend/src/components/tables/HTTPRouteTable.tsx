"use client";

import React from "react";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableCell,
} from "@/components/ui/table";
import Badge from "@/components/ui/badge/Badge";
import { HTTPRouteDetail } from "@/lib/api";

interface HTTPRouteTableProps {
  routes: HTTPRouteDetail[];
  loading: boolean;
  onDelete: (namespace: string, name: string) => void;
}

export default function HTTPRouteTable({ routes, loading, onDelete }: HTTPRouteTableProps) {
  const headers = (
    <TableHeader>
      <TableRow className="border-b border-gray-100 dark:border-gray-800">
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          Name
        </TableCell>
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          Hostnames
        </TableCell>
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          Parent Gateway
        </TableCell>
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          Actions
        </TableCell>
      </TableRow>
    </TableHeader>
  );

  return (
    <div className="overflow-hidden rounded-xl border border-gray-200 dark:border-gray-800">
      <div className="max-w-full overflow-x-auto">
        <Table>
          {headers}
          <TableBody>
            {loading ? (
              <tr>
                <td className="px-5 py-12 text-center text-sm text-gray-500 dark:text-gray-400" colSpan={4}>
                  Loading...
                </td>
              </tr>
            ) : routes.length === 0 ? (
              <tr>
                <td className="px-5 py-12 text-center text-sm text-gray-500 dark:text-gray-400" colSpan={4}>
                  No HTTPRoutes configured for this service
                </td>
              </tr>
            ) : (
              routes.map((route) => (
                <TableRow
                  key={`${route.namespace}/${route.name}`}
                  className="border-b border-gray-100 dark:border-gray-800 last:border-b-0"
                >
                  <TableCell className="px-5 py-4 text-sm font-medium text-gray-800 dark:text-white/90">
                    {route.name}
                  </TableCell>
                  <TableCell className="px-5 py-4">
                    <div className="flex flex-wrap gap-1">
                      {route.hostnames.map((h) => (
                        <Badge key={h} variant="light" color="primary" size="sm">
                          {h}
                        </Badge>
                      ))}
                    </div>
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm text-gray-600 dark:text-gray-300">
                    {route.parentRefs.map((ref) => (
                      <span key={ref.name} className="block">
                        {ref.namespace ? `${ref.namespace}/` : ""}
                        {ref.name}
                      </span>
                    ))}
                  </TableCell>
                  <TableCell className="px-5 py-4">
                    <button
                      onClick={() => onDelete(route.namespace, route.name)}
                      className="text-sm text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
                    >
                      Delete
                    </button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
