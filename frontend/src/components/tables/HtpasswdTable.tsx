"use client";

import React from "react";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableCell,
} from "@/components/ui/table";
import type { HtpasswdSecretSummary } from "@/lib/api";

interface HtpasswdTableProps {
  secrets: HtpasswdSecretSummary[];
  loading: boolean;
  onSelect: (name: string) => void;
  selectedName: string | undefined;
}

export default function HtpasswdTable({ secrets, loading, onSelect, selectedName }: HtpasswdTableProps) {
  const formatDate = (isoString: string): string => {
    const date = new Date(isoString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return (
    <div className="overflow-hidden rounded-xl border border-gray-200 dark:border-gray-800">
      <div className="max-w-full overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow className="border-b border-gray-100 dark:border-gray-800">
              <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Name
              </TableCell>
              <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                User Count
              </TableCell>
              <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Created At
              </TableCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <tr>
                <td className="px-5 py-12 text-center text-sm text-gray-500 dark:text-gray-400" colSpan={3}>
                  Loading...
                </td>
              </tr>
            ) : secrets.length === 0 ? (
              <tr>
                <td className="px-5 py-12 text-center text-sm text-gray-500 dark:text-gray-400" colSpan={3}>
                  No htpasswd secrets found
                </td>
              </tr>
            ) : (
              secrets.map((secret) => (
                <tr
                  key={secret.name}
                  className={`border-b border-gray-100 dark:border-gray-800 last:border-b-0 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800/50 ${
                    selectedName === secret.name ? "bg-brand-50 dark:bg-brand-900/20" : ""
                  }`}
                  onClick={() => onSelect(secret.name)}
                >
                  <TableCell className="px-5 py-4 text-sm font-medium text-gray-800 dark:text-white/90">
                    {secret.name}
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm text-gray-600 dark:text-gray-300">
                    {secret.userCount}
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm text-gray-600 dark:text-gray-300">
                    {formatDate(secret.createdAt)}
                  </TableCell>
                </tr>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
