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
import type { ServiceOverview } from "@/lib/api";

interface ServiceTableProps {
  services: ServiceOverview[];
  loading: boolean;
}

export default function ServiceTable({ services, loading }: ServiceTableProps) {
  const headers = (
    <TableHeader>
      <TableRow className="border-b border-gray-100 dark:border-gray-800">
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          Name
        </TableCell>
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          Namespace
        </TableCell>
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          Cluster IP
        </TableCell>
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          Ports
        </TableCell>
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          HTTPRoute
        </TableCell>
        <TableCell isHeader className="px-5 py-3 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          SecurityPolicy
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
                <td className="px-5 py-12 text-center text-sm text-gray-500 dark:text-gray-400" colSpan={6}>
                  Loading...
                </td>
              </tr>
            ) : services.length === 0 ? (
              <tr>
                <td className="px-5 py-12 text-center text-sm text-gray-500 dark:text-gray-400" colSpan={6}>
                  No services found
                </td>
              </tr>
            ) : (
              services.map((svc) => (
                <TableRow
                  key={`${svc.namespace}/${svc.name}`}
                  className="border-b border-gray-100 dark:border-gray-800 last:border-b-0"
                >
                  <TableCell className="px-5 py-4 text-sm font-medium text-gray-800 dark:text-white/90">
                    {svc.name}
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm text-gray-600 dark:text-gray-300">
                    {svc.namespace}
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm font-mono text-gray-600 dark:text-gray-300">
                    {svc.clusterIP}
                  </TableCell>
                  <TableCell className="px-5 py-4 text-sm text-gray-600 dark:text-gray-300">
                    {svc.ports.map((p) => `${p.port}/${p.targetPort}`).join(", ")}
                  </TableCell>
                  <TableCell className="px-5 py-4">
                    {svc.httpRoute ? (
                      <div className="flex flex-col gap-1">
                        {svc.httpRoute.hostnames.map((h) => (
                          <Badge key={h} variant="light" color="primary" size="sm">
                            {h}
                          </Badge>
                        ))}
                      </div>
                    ) : (
                      <Badge variant="light" color="light" size="sm">
                        None
                      </Badge>
                    )}
                  </TableCell>
                  <TableCell className="px-5 py-4">
                    {svc.securityPolicy ? (
                      <div className="flex flex-wrap gap-1">
                        {svc.securityPolicy.hasBasicAuth && (
                          <Badge variant="light" color="success" size="sm">
                            BasicAuth
                          </Badge>
                        )}
                        {svc.securityPolicy.hasTLS && (
                          <Badge variant="light" color="info" size="sm">
                            TLS
                          </Badge>
                        )}
                        {!svc.securityPolicy.hasBasicAuth && !svc.securityPolicy.hasTLS && (
                          <Badge variant="light" color="warning" size="sm">
                            Configured
                          </Badge>
                        )}
                      </div>
                    ) : (
                      <Badge variant="light" color="light" size="sm">
                        None
                      </Badge>
                    )}
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
